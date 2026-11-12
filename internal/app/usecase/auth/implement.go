package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"go-react-boilerplate/internal/app/port"
	"go-react-boilerplate/internal/entity"
)

const (
	challengeTTL         = 5 * time.Minute
	passwordConfirmTTL   = 15 * time.Minute
	sessionUpdateAge     = 24 * time.Hour
	twoFactorMaxAttempts = 5
	twoFactorLockout     = 5 * time.Minute
)

type Config struct {
	AppURL     string
	SessionTTL time.Duration
	Secret     []byte
}

type useCase struct {
	users         port.UserStore
	roles         port.RoleStore
	settings      port.SettingStore
	sessions      port.SessionStore
	verifications port.VerificationStore
	twoFactors    port.TwoFactorStore
	hasher        port.PasswordHasher
	tokens        port.TokenService
	totp          port.TOTPService
	cipher        port.SecretCipher
	mailer        port.Mailer
	limiter       port.RateLimiter
	clock         port.Clock
	cfg           Config
	dummyHash     string
}

func New(
	users port.UserStore,
	roles port.RoleStore,
	settings port.SettingStore,
	sessions port.SessionStore,
	verifications port.VerificationStore,
	twoFactors port.TwoFactorStore,
	hasher port.PasswordHasher,
	tokens port.TokenService,
	totp port.TOTPService,
	cipher port.SecretCipher,
	mailer port.Mailer,
	limiter port.RateLimiter,
	clock port.Clock,
	cfg Config,
) UseCase {
	if cfg.SessionTTL <= 0 {
		cfg.SessionTTL = 7 * 24 * time.Hour
	}
	dummy, err := hasher.Hash("dummy-password-for-timing")
	if err != nil {
		dummy = "invalid"
	}
	return &useCase{
		users: users, roles: roles, settings: settings, sessions: sessions,
		verifications: verifications, twoFactors: twoFactors, hasher: hasher,
		tokens: tokens, totp: totp, cipher: cipher, mailer: mailer,
		limiter: limiter, clock: clock, cfg: cfg, dummyHash: dummy,
	}
}

func (u *useCase) Register(ctx context.Context, input RegisterInput) (RegisterResult, error) {
	if !u.limiter.Allow("register:" + input.IP) {
		return RegisterResult{}, entity.ErrTooManyRequests
	}

	name := strings.TrimSpace(input.Name)
	email := entity.NormalizeEmail(input.Email)

	fields := map[string]string{}
	if !entity.ValidateName(name) {
		fields["name"] = "validation.nameLength"
	}
	if !entity.ValidEmail(email) {
		fields["email"] = "validation.email"
	}
	if !entity.ValidPassword(input.Password) {
		fields["password"] = "validation.passwordLength"
	}
	if len(fields) > 0 {
		return RegisterResult{}, entity.ValidationError(fields)
	}

	if _, err := u.users.GetByEmail(ctx, email); err == nil {
		return RegisterResult{}, entity.ErrEmailInUse
	} else if !errors.Is(err, entity.ErrNotFound) {
		return RegisterResult{}, err
	}

	role, err := u.roles.GetByName(ctx, entity.RoleUser)
	if err != nil {
		return RegisterResult{}, err
	}

	approval, err := u.registrationApproval(ctx)
	if err != nil {
		return RegisterResult{}, err
	}

	status := entity.StatusActive
	if approval {
		status = entity.StatusPending
	}

	hash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return RegisterResult{}, err
	}

	now := u.clock.Now().UTC()
	user := &entity.User{
		ID: entity.NewUUID(), Name: name, Email: email, RoleID: role.ID,
		Status: status, CreatedAt: now, UpdatedAt: now,
	}
	if err := u.users.CreateWithPassword(ctx, user, hash); err != nil {
		return RegisterResult{}, err
	}

	sent := u.sendVerification(ctx, user)
	return RegisterResult{Status: status, EmailSent: sent}, nil
}

func (u *useCase) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	if !u.limiter.Allow("login:" + input.IP) {
		return LoginResult{}, entity.ErrTooManyRequests
	}

	email := entity.NormalizeEmail(input.Email)
	user, err := u.users.GetByEmail(ctx, email)
	if err != nil {
		// Verify against a dummy hash so an unknown address costs the same as
		// a wrong password. Do not reveal which one failed.
		u.hasher.Verify(u.dummyHash, input.Password)
		if errors.Is(err, entity.ErrNotFound) {
			return LoginResult{}, entity.ErrInvalidCredentials
		}
		return LoginResult{}, err
	}

	hash, err := u.users.GetPasswordHash(ctx, user.ID)
	if err != nil {
		u.hasher.Verify(u.dummyHash, input.Password)
		return LoginResult{}, entity.ErrInvalidCredentials
	}
	if !u.hasher.Verify(hash, input.Password) {
		return LoginResult{}, entity.ErrInvalidCredentials
	}

	switch user.Status {
	case entity.StatusPending:
		return LoginResult{}, entity.ErrAccountPending
	case entity.StatusSuspended:
		return LoginResult{}, entity.ErrAccountSuspended
	case entity.StatusActive:
	default:
		return LoginResult{}, entity.ErrUnauthorized
	}

	twoFactor, err := u.twoFactors.GetByUserID(ctx, user.ID)
	if err == nil && twoFactor.Verified {
		challenge, err := u.signValue(user.ID, challengeTTL)
		if err != nil {
			return LoginResult{}, err
		}
		return LoginResult{TwoFactorRequired: true, Challenge: challenge}, nil
	}
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return LoginResult{}, err
	}

	return u.createSession(ctx, user, input.IP, input.UserAgent)
}

func (u *useCase) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	err := u.sessions.DeleteByTokenHash(ctx, u.tokens.HashToken(token))
	if errors.Is(err, entity.ErrNotFound) {
		return nil
	}
	return err
}

func (u *useCase) Session(ctx context.Context, token string) (*entity.SessionContext, error) {
	if token == "" {
		return nil, entity.ErrUnauthorized
	}
	session, err := u.sessions.GetContext(ctx, u.tokens.HashToken(token))
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, entity.ErrUnauthorized
		}
		return nil, err
	}
	now := u.clock.Now().UTC()
	if !session.Session.Active(now) {
		_ = u.sessions.DeleteByTokenHash(ctx, session.Session.TokenHash)
		return nil, entity.ErrUnauthorized
	}
	if session.User.Status != entity.StatusActive {
		_ = u.sessions.DeleteByTokenHash(ctx, session.Session.TokenHash)
		return nil, entity.ErrUnauthorized
	}
	// Sliding expiry: touch at most once a day so we do not write on every
	// request.
	if now.Sub(session.Session.UpdatedAt) > sessionUpdateAge {
		expires := now.Add(u.cfg.SessionTTL)
		if err := u.sessions.UpdateExpiry(ctx, session.Session.ID, expires); err == nil {
			session.Session.ExpiresAt = expires
			session.Session.UpdatedAt = now
		}
	}
	return session, nil
}

func (u *useCase) VerifyEmail(ctx context.Context, token string) error {
	if token == "" {
		return entity.ErrInvalidToken
	}
	verification, err := u.verifications.Consume(ctx, u.tokens.HashToken(token), entity.VerifyEmail, u.clock.Now().UTC())
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrInvalidToken
		}
		return err
	}
	user, err := u.users.GetByID(ctx, verification.Identifier)
	if err != nil {
		return entity.ErrInvalidToken
	}
	user.EmailVerified = true
	return u.users.UpdateWithPassword(ctx, user, "")
}

func (u *useCase) RequestEmailVerification(ctx context.Context, userID, ip string) error {
	if ip != "" && !u.limiter.Allow("verify-request:"+ip) {
		return entity.ErrTooManyRequests
	}
	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		return entity.ErrUserNotFound
	}
	u.sendVerification(ctx, user)
	return nil
}

func (u *useCase) ForgotPassword(ctx context.Context, email, ip string) error {
	if !u.limiter.Allow("forgot:" + ip) {
		return entity.ErrTooManyRequests
	}
	user, err := u.users.GetByEmail(ctx, entity.NormalizeEmail(email))
	if err != nil {
		// Always report success: a reset form must not enumerate accounts.
		return nil
	}
	token, err := u.tokens.NewToken()
	if err != nil {
		return err
	}
	_ = u.verifications.DeleteForIdentifier(ctx, user.ID, entity.ResetPassword)
	now := u.clock.Now().UTC()
	if err := u.verifications.Create(ctx, &entity.Verification{
		ID: entity.NewUUID(), Identifier: user.ID, Purpose: entity.ResetPassword,
		TokenHash: u.tokens.HashToken(token), ExpiresAt: now.Add(time.Hour), CreatedAt: now,
	}); err != nil {
		return err
	}
	link := u.link("/reset-password", token)
	_ = u.mailer.SendPasswordReset(ctx, user.Email, user.Name, link)
	return nil
}

func (u *useCase) ResetPassword(ctx context.Context, token, password string) error {
	if !entity.ValidPassword(password) {
		return entity.ValidationError(map[string]string{"password": "validation.passwordLength"})
	}
	verification, err := u.verifications.Consume(ctx, u.tokens.HashToken(token), entity.ResetPassword, u.clock.Now().UTC())
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrInvalidToken
		}
		return err
	}
	user, err := u.users.GetByID(ctx, verification.Identifier)
	if err != nil {
		return entity.ErrInvalidToken
	}
	hash, err := u.hasher.Hash(password)
	if err != nil {
		return err
	}
	if err := u.users.UpdateWithPassword(ctx, user, hash); err != nil {
		return err
	}
	return u.sessions.DeleteByUserID(ctx, user.ID)
}

func (u *useCase) VerifyTwoFactor(ctx context.Context, challenge, code string, backup bool, ip, userAgent string) (LoginResult, error) {
	if !u.limiter.Allow("2fa:" + ip) {
		return LoginResult{}, entity.ErrTooManyRequests
	}
	userID, err := u.verifyValue(challenge)
	if err != nil {
		return LoginResult{}, entity.ErrUnauthorized
	}
	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		return LoginResult{}, entity.ErrUnauthorized
	}
	switch user.Status {
	case entity.StatusPending:
		return LoginResult{}, entity.ErrAccountPending
	case entity.StatusSuspended:
		return LoginResult{}, entity.ErrAccountSuspended
	case entity.StatusActive:
	default:
		return LoginResult{}, entity.ErrUnauthorized
	}

	record, err := u.twoFactors.GetByUserID(ctx, user.ID)
	if err != nil || !record.Verified {
		return LoginResult{}, entity.ErrTwoFactorNotEnabled
	}
	now := u.clock.Now().UTC()
	if record.LockedUntil != nil && now.Before(*record.LockedUntil) {
		return LoginResult{}, entity.ErrTooManyRequests
	}

	ok := false
	if backup {
		hashed := u.tokens.HashBackupCode(strings.TrimSpace(code))
		for index, candidate := range record.BackupCodeHashes {
			if subtle.ConstantTimeCompare([]byte(candidate), []byte(hashed)) == 1 {
				record.BackupCodeHashes = append(record.BackupCodeHashes[:index], record.BackupCodeHashes[index+1:]...)
				ok = true
				break
			}
		}
	} else {
		secret, err := u.cipher.Decrypt(record.Secret)
		if err != nil {
			return LoginResult{}, entity.ErrUnexpected
		}
		ok = u.totp.Validate(secret, strings.TrimSpace(code))
	}

	if !ok {
		record.FailedVerificationCount++
		if record.FailedVerificationCount >= twoFactorMaxAttempts {
			locked := now.Add(twoFactorLockout)
			record.LockedUntil = &locked
			record.FailedVerificationCount = 0
		}
		_ = u.twoFactors.Save(ctx, record)
		return LoginResult{}, entity.ErrInvalidCode
	}

	record.FailedVerificationCount = 0
	record.LockedUntil = nil
	if err := u.twoFactors.Save(ctx, record); err != nil {
		return LoginResult{}, err
	}

	return u.createSession(ctx, user, ip, userAgent)
}

func (u *useCase) ConfirmPassword(ctx context.Context, session *entity.SessionContext, password string) (string, error) {
	hash, err := u.users.GetPasswordHash(ctx, session.User.ID)
	if err != nil {
		return "", entity.ErrIncorrectPassword
	}
	if !u.hasher.Verify(hash, password) {
		return "", entity.ErrIncorrectPassword
	}
	return u.signValue(session.User.ID, passwordConfirmTTL)
}

func (u *useCase) HasRecentPasswordConfirmation(userID, cookieValue string, now time.Time) bool {
	parts := strings.Split(cookieValue, ".")
	if len(parts) != 3 {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(parts[0]), []byte(userID)) != 1 {
		return false
	}
	issued, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false
	}
	if now.Sub(time.Unix(issued, 0)) > passwordConfirmTTL {
		return false
	}
	expected := u.signRaw(parts[0] + "." + parts[1])
	return subtle.ConstantTimeCompare([]byte(expected), []byte(parts[2])) == 1
}

func (u *useCase) createSession(ctx context.Context, user *entity.User, ip, userAgent string) (LoginResult, error) {
	raw, err := u.tokens.NewToken()
	if err != nil {
		return LoginResult{}, err
	}
	now := u.clock.Now().UTC()
	session := &entity.Session{
		ID: entity.NewUUID(), UserID: user.ID, TokenHash: u.tokens.HashToken(raw),
		ExpiresAt: now.Add(u.cfg.SessionTTL), CreatedAt: now, UpdatedAt: now,
		IPAddress: ip, UserAgent: userAgent,
	}
	if err := u.sessions.Create(ctx, session); err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		Token: raw,
		Context: &entity.SessionContext{
			Session: *session, User: *user, Permissions: user.Permissions,
		},
	}, nil
}

func (u *useCase) registrationApproval(ctx context.Context) (bool, error) {
	value, err := u.settings.Get(ctx, entity.SettingRegistrationApproval)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return entity.RegistrationApprovalEnabled(value), nil
}

func (u *useCase) sendVerification(ctx context.Context, user *entity.User) bool {
	token, err := u.tokens.NewToken()
	if err != nil {
		return false
	}
	_ = u.verifications.DeleteForIdentifier(ctx, user.ID, entity.VerifyEmail)
	now := u.clock.Now().UTC()
	if err := u.verifications.Create(ctx, &entity.Verification{
		ID: entity.NewUUID(), Identifier: user.ID, Purpose: entity.VerifyEmail,
		TokenHash: u.tokens.HashToken(token), ExpiresAt: now.Add(time.Hour), CreatedAt: now,
	}); err != nil {
		return false
	}
	link := u.link("/verify-email", token)
	return u.mailer.SendVerification(ctx, user.Email, user.Name, link) == nil
}

func (u *useCase) link(path, token string) string {
	base := strings.TrimRight(u.cfg.AppURL, "/")
	return base + path + "?token=" + token
}

func (u *useCase) signValue(value string, ttl time.Duration) (string, error) {
	expires := strconv.FormatInt(u.clock.Now().UTC().Add(ttl).Unix(), 10)
	payload := value + "." + expires
	return payload + "." + u.signRaw(payload), nil
}

func (u *useCase) verifyValue(value string) (string, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return "", entity.ErrUnauthorized
	}
	expected := u.signRaw(parts[0] + "." + parts[1])
	if subtle.ConstantTimeCompare([]byte(expected), []byte(parts[2])) != 1 {
		return "", entity.ErrUnauthorized
	}
	expires, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || u.clock.Now().UTC().After(time.Unix(expires, 0)) {
		return "", entity.ErrUnauthorized
	}
	return parts[0], nil
}

func (u *useCase) signRaw(payload string) string {
	mac := hmac.New(sha256.New, u.cfg.Secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
