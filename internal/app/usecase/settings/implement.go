package settings

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"go-react-boilerplate/internal/app/port"
	"go-react-boilerplate/internal/entity"
)

func (u *useCase) UpdateProfile(ctx context.Context, session *entity.SessionContext, input ProfileInput) error {
	fields := map[string]string{}
	name := strings.TrimSpace(input.Name)
	email := entity.NormalizeEmail(input.Email)
	if !entity.ValidateName(name) {
		fields["name"] = "validation.nameLength"
	}
	if !entity.ValidEmail(email) {
		fields["email"] = "validation.email"
	}
	if len([]rune(strings.TrimSpace(input.Phone))) > 40 {
		fields["phone"] = "validation.max"
	}
	if len([]rune(strings.TrimSpace(input.JobTitle))) > 120 {
		fields["jobTitle"] = "validation.max"
	}
	if len([]rune(strings.TrimSpace(input.Department))) > 120 {
		fields["department"] = "validation.max"
	}
	if len(fields) > 0 {
		return entity.ValidationError(fields)
	}

	user, err := u.users.GetByID(ctx, session.User.ID)
	if err != nil {
		return entity.ErrUserNotFound
	}
	if email != user.Email {
		existing, err := u.users.GetByEmail(ctx, email)
		if err == nil && existing.ID != user.ID {
			return entity.ErrEmailInUse
		}
		if err != nil && !errors.Is(err, entity.ErrNotFound) {
			return err
		}
	}

	emailChanged := email != user.Email
	user.Name = name
	user.Email = email
	if emailChanged {
		user.EmailVerified = false
	}
	user.PhoneNumber = strings.TrimSpace(input.Phone)
	user.JobTitle = strings.TrimSpace(input.JobTitle)
	user.Department = strings.TrimSpace(input.Department)
	user.UpdatedAt = u.clock.Now().UTC()

	if err := u.users.UpdateWithPassword(ctx, user, ""); err != nil {
		if errors.Is(err, port.ErrUniqueViolation) {
			return entity.ErrEmailInUse
		}
		return err
	}
	if emailChanged {
		// A mail server failure must not roll back a profile edit. The
		// verifier logs the fallback link.
		_ = u.verifier.RequestEmailVerification(ctx, user.ID, "")
	}
	return nil
}

func (u *useCase) ChangePassword(ctx context.Context, session *entity.SessionContext, input PasswordInput) error {
	if !entity.ValidPassword(input.Next) {
		return entity.ValidationError(map[string]string{"newPassword": "validation.passwordLength"})
	}
	hash, err := u.users.GetPasswordHash(ctx, session.User.ID)
	if err != nil {
		return entity.ErrIncorrectPassword
	}
	if !u.hasher.Verify(hash, input.Current) {
		return entity.ErrIncorrectPassword
	}
	next, err := u.hasher.Hash(input.Next)
	if err != nil {
		return err
	}
	user, err := u.users.GetByID(ctx, session.User.ID)
	if err != nil {
		return entity.ErrUserNotFound
	}
	user.UpdatedAt = u.clock.Now().UTC()
	if err := u.users.UpdateWithPassword(ctx, user, next); err != nil {
		return err
	}
	// Revoke every other session. The current cookie keeps working because
	// its row is the one we keep.
	return u.sessions.DeleteOthers(ctx, user.ID, session.Session.ID)
}

func (u *useCase) AccountOverview(ctx context.Context, session *entity.SessionContext, confirmed bool) (AccountOverview, error) {
	canDelete := true
	if session.User.Role != nil && session.User.Role.Name == entity.RoleSuperadmin {
		count, err := u.users.CountSuperadmins(ctx)
		if err != nil {
			return AccountOverview{}, err
		}
		canDelete = count > 1
	}
	return AccountOverview{CanDelete: canDelete, Confirmed: confirmed}, nil
}

func (u *useCase) DeleteAccount(ctx context.Context, session *entity.SessionContext, confirmed bool) error {
	if !confirmed {
		return entity.ErrConfirmPasswordRequired
	}
	if session.User.Role != nil && session.User.Role.Name == entity.RoleSuperadmin {
		count, err := u.users.CountSuperadmins(ctx)
		if err != nil {
			return err
		}
		if count <= 1 {
			return entity.ErrLastSuperadmin
		}
	}
	return u.users.Delete(ctx, session.User.ID)
}

func (u *useCase) EnableTwoFactor(ctx context.Context, session *entity.SessionContext, password string) (TwoFactorSetup, error) {
	hash, err := u.users.GetPasswordHash(ctx, session.User.ID)
	if err != nil || !u.hasher.Verify(hash, password) {
		return TwoFactorSetup{}, entity.ErrIncorrectPassword
	}
	if existing, err := u.twoFactor.GetByUserID(ctx, session.User.ID); err == nil && existing.Verified {
		return TwoFactorSetup{}, entity.NewError("settings.security.alreadyEnabled")
	}

	secret, uri, err := u.totp.Generate(u.appName, session.User.Email)
	if err != nil {
		return TwoFactorSetup{}, err
	}
	encrypted, err := u.cipher.Encrypt(secret)
	if err != nil {
		return TwoFactorSetup{}, err
	}
	codes := make([]string, 0, 10)
	hashes := make([]string, 0, 10)
	for range 10 {
		code, err := u.tokens.NewBackupCode()
		if err != nil {
			return TwoFactorSetup{}, err
		}
		codes = append(codes, code)
		hashes = append(hashes, u.tokens.HashBackupCode(code))
	}
	now := u.clock.Now().UTC()
	record := &entity.TwoFactor{
		ID: entity.NewUUID(), UserID: session.User.ID, Secret: encrypted,
		BackupCodeHashes: hashes, Verified: false, CreatedAt: now, UpdatedAt: now,
	}
	if err := u.twoFactor.Save(ctx, record); err != nil {
		return TwoFactorSetup{}, err
	}
	qr := ""
	if png, err := u.totp.QRCode(uri); err == nil {
		qr = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	}
	return TwoFactorSetup{Secret: secret, URI: uri, QRDataURL: qr, BackupCodes: codes}, nil
}

func (u *useCase) ConfirmTwoFactor(ctx context.Context, session *entity.SessionContext, code string) error {
	record, err := u.twoFactor.GetByUserID(ctx, session.User.ID)
	if err != nil {
		return entity.ErrTwoFactorNotEnabled
	}
	secret, err := u.cipher.Decrypt(record.Secret)
	if err != nil {
		return entity.ErrUnexpected
	}
	if !u.totp.Validate(secret, strings.TrimSpace(code)) {
		return entity.ErrInvalidCode
	}
	record.Verified = true
	record.UpdatedAt = u.clock.Now().UTC()
	if err := u.twoFactor.Save(ctx, record); err != nil {
		return err
	}
	user, err := u.users.GetByID(ctx, session.User.ID)
	if err != nil {
		return entity.ErrUserNotFound
	}
	user.TwoFactorEnabled = true
	user.UpdatedAt = u.clock.Now().UTC()
	return u.users.UpdateWithPassword(ctx, user, "")
}

func (u *useCase) DisableTwoFactor(ctx context.Context, session *entity.SessionContext, password string) error {
	hash, err := u.users.GetPasswordHash(ctx, session.User.ID)
	if err != nil || !u.hasher.Verify(hash, password) {
		return entity.ErrIncorrectPassword
	}
	if err := u.twoFactor.DeleteByUserID(ctx, session.User.ID); err != nil {
		return err
	}
	user, err := u.users.GetByID(ctx, session.User.ID)
	if err != nil {
		return entity.ErrUserNotFound
	}
	user.TwoFactorEnabled = false
	user.UpdatedAt = u.clock.Now().UTC()
	return u.users.UpdateWithPassword(ctx, user, "")
}
