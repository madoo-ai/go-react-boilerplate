// Package auth implements registration, login, email verification, password
// reset, TOTP two-factor verification and the confirm-password gate.
package auth

import (
	"context"
	"time"

	"go-react-boilerplate/internal/entity"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
	IP       string
}

type RegisterResult struct {
	Status    entity.UserStatus
	EmailSent bool
}

type LoginInput struct {
	Email      string
	Password   string
	IP         string
	UserAgent  string
	RememberMe bool
}

type LoginResult struct {
	// Token is the raw session token. It is empty when two-factor
	// verification is required.
	Token string
	// Context is the authenticated user when Token is non-empty.
	Context *entity.SessionContext
	// TwoFactorRequired tells the client to finish the flow on
	// /two-factor. Challenge is a short-lived signed cookie value.
	TwoFactorRequired bool
	Challenge         string
}

type UseCase interface {
	Register(ctx context.Context, input RegisterInput) (RegisterResult, error)
	Login(ctx context.Context, input LoginInput) (LoginResult, error)
	Logout(ctx context.Context, token string) error
	Session(ctx context.Context, token string) (*entity.SessionContext, error)
	VerifyEmail(ctx context.Context, token string) error
	RequestEmailVerification(ctx context.Context, userID, ip string) error
	ForgotPassword(ctx context.Context, email, ip string) error
	ResetPassword(ctx context.Context, token, password string) error
	VerifyTwoFactor(ctx context.Context, challenge, code string, backup bool, ip, userAgent string) (LoginResult, error)
	// ConfirmPassword verifies the current password and returns the value for
	// the password-confirmation cookie.
	ConfirmPassword(ctx context.Context, session *entity.SessionContext, password string) (string, error)
	// HasRecentPasswordConfirmation checks a previously issued cookie value.
	HasRecentPasswordConfirmation(userID, cookieValue string, now time.Time) bool
}
