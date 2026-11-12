package port

import (
	"context"
	"errors"
	"time"

	"go-react-boilerplate/internal/entity"
)

// Clock is injected so usecases and tests share one notion of "now".
type Clock interface {
	Now() time.Time
}

// ErrUniqueViolation is reported by stores when a unique constraint rejects a
// write (for example a duplicate email). Adapters translate their driver's
// error into this sentinel so usecases stay portable.
var ErrUniqueViolation = errors.New("unique violation")

// PasswordHasher hides the password KDF. The production adapter is Argon2id.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hash, password string) bool
}

// TokenService mints opaque tokens and hashes them for storage. Raw tokens
// only ever leave the process in emails or cookies.
type TokenService interface {
	NewToken() (string, error)
	HashToken(token string) string
	NewBackupCode() (string, error)
	HashBackupCode(code string) string
}

// SecretCipher encrypts TOTP secrets at rest.
type SecretCipher interface {
	Encrypt(plain string) (string, error)
	Decrypt(encoded string) (string, error)
}

// TOTPService generates and verifies time-based one-time passwords.
type TOTPService interface {
	Generate(issuer, account string) (secret, uri string, err error)
	Validate(secret, code string) bool
	QRCode(uri string) ([]byte, error)
}

// RateLimiter is a per-key sliding window. Keys combine an IP and an action.
type RateLimiter interface {
	Allow(key string) bool
}

// Mailer delivers the two auth emails. Implementations must not block the
// request path when the local mail server is missing.
type Mailer interface {
	SendVerification(ctx context.Context, to, name, url string) error
	SendPasswordReset(ctx context.Context, to, name, url string) error
}

// UserStore persists accounts and their credentials.
type UserStore interface {
	CreateWithPassword(ctx context.Context, user *entity.User, passwordHash string) error
	UpdateWithPassword(ctx context.Context, user *entity.User, passwordHash string) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	List(ctx context.Context, filters entity.UserFilters) ([]entity.User, int, error)
	CountSuperadmins(ctx context.Context) (int, error)
	GetPasswordHash(ctx context.Context, userID string) (string, error)
}

// RoleStore persists roles and the permission matrix.
type RoleStore interface {
	List(ctx context.Context) ([]entity.Role, error)
	GetByID(ctx context.Context, id string) (*entity.Role, error)
	GetByName(ctx context.Context, name string) (*entity.Role, error)
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id string) error
	ListPermissions(ctx context.Context) ([]entity.PermissionInfo, error)
}

// SettingStore is a small key/value store.
type SettingStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}

// SessionStore persists server-side sessions. GetContext loads the session,
// user, role and permissions in one call.
type SessionStore interface {
	Create(ctx context.Context, session *entity.Session) error
	GetContext(ctx context.Context, tokenHash string) (*entity.SessionContext, error)
	UpdateExpiry(ctx context.Context, id string, expiresAt time.Time) error
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteOthers(ctx context.Context, userID, keepSessionID string) error
	DeleteExpired(ctx context.Context, now time.Time) error
}

// VerificationStore holds one-shot email tokens.
type VerificationStore interface {
	Create(ctx context.Context, verification *entity.Verification) error
	Consume(ctx context.Context, tokenHash string, purpose entity.VerificationPurpose, now time.Time) (*entity.Verification, error)
	DeleteForIdentifier(ctx context.Context, identifier string, purpose entity.VerificationPurpose) error
}

// TwoFactorStore persists TOTP state.
type TwoFactorStore interface {
	GetByUserID(ctx context.Context, userID string) (*entity.TwoFactor, error)
	Save(ctx context.Context, twoFactor *entity.TwoFactor) error
	DeleteByUserID(ctx context.Context, userID string) error
}

// DashboardStore powers the admin overview. It stays separate from UserStore
// so the dashboard queries can evolve without widening the account port.
type DashboardStore interface {
	Stats(ctx context.Context) (total, active, pending int, err error)
	RoleCounts(ctx context.Context) ([]entity.RoleDatum, error)
	SignupTimes(ctx context.Context, since time.Time) ([]time.Time, error)
	RecentUsers(ctx context.Context, limit int) ([]entity.User, error)
}
