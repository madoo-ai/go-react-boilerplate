package entity

import "time"

// Session is a server-side login record. Only a hash of the cookie token is
// stored, so a database leak does not hand out live sessions.
type Session struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	IPAddress string
	UserAgent string
}

func (s Session) Active(now time.Time) bool {
	return now.Before(s.ExpiresAt)
}

// SessionContext is the current session joined with its user, role and
// permissions. It is what every authenticated handler receives.
type SessionContext struct {
	Session     Session
	User        User
	Permissions []Permission
}

func (s SessionContext) Actor() Actor {
	roleName := ""
	if s.User.Role != nil {
		roleName = s.User.Role.Name
	}
	return Actor{
		ID: s.User.ID, Name: s.User.Name, Email: s.User.Email,
		RoleName: roleName, Permissions: s.Permissions,
	}
}

// VerificationPurpose separates the two emailed token flows that share the
// verifications table.
type VerificationPurpose string

const (
	VerifyEmail   VerificationPurpose = "email_verification"
	ResetPassword VerificationPurpose = "password_reset"
)

type Verification struct {
	ID         string
	Identifier string
	Purpose    VerificationPurpose
	TokenHash  string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}

// TwoFactor is the stored TOTP state. Secret is encrypted at rest by the
// security adapter. BackupCodeHashes are SHA-256 hex digests.
type TwoFactor struct {
	ID                      string
	UserID                  string
	Secret                  string
	BackupCodeHashes        []string
	Verified                bool
	FailedVerificationCount int
	LockedUntil             *time.Time
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
