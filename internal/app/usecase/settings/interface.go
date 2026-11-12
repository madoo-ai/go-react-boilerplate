// Package settings implements self-service profile, password, two-factor and
// account-deletion use cases.
package settings

import (
	"context"

	"go-react-boilerplate/internal/app/port"
	"go-react-boilerplate/internal/entity"
)

type EmailVerifier interface {
	RequestEmailVerification(ctx context.Context, userID, ip string) error
}

type ProfileInput struct {
	Name       string
	Email      string
	Phone      string
	JobTitle   string
	Department string
}

type PasswordInput struct {
	Current string
	Next    string
}

type AccountOverview struct {
	CanDelete bool
	Confirmed bool
}

type TwoFactorSetup struct {
	Secret      string
	URI         string
	QRDataURL   string
	BackupCodes []string
}

type UseCase interface {
	UpdateProfile(ctx context.Context, session *entity.SessionContext, input ProfileInput) error
	ChangePassword(ctx context.Context, session *entity.SessionContext, input PasswordInput) error
	AccountOverview(ctx context.Context, session *entity.SessionContext, confirmed bool) (AccountOverview, error)
	DeleteAccount(ctx context.Context, session *entity.SessionContext, confirmed bool) error
	EnableTwoFactor(ctx context.Context, session *entity.SessionContext, password string) (TwoFactorSetup, error)
	ConfirmTwoFactor(ctx context.Context, session *entity.SessionContext, code string) error
	DisableTwoFactor(ctx context.Context, session *entity.SessionContext, password string) error
}

type useCase struct {
	users     port.UserStore
	sessions  port.SessionStore
	twoFactor port.TwoFactorStore
	hasher    port.PasswordHasher
	tokens    port.TokenService
	totp      port.TOTPService
	cipher    port.SecretCipher
	verifier  EmailVerifier
	clock     port.Clock
	appName   string
}

func New(
	users port.UserStore,
	sessions port.SessionStore,
	twoFactor port.TwoFactorStore,
	hasher port.PasswordHasher,
	tokens port.TokenService,
	totp port.TOTPService,
	cipher port.SecretCipher,
	verifier EmailVerifier,
	clock port.Clock,
	appName string,
) UseCase {
	return &useCase{
		users: users, sessions: sessions, twoFactor: twoFactor, hasher: hasher,
		tokens: tokens, totp: totp, cipher: cipher, verifier: verifier, clock: clock,
		appName: appName,
	}
}
