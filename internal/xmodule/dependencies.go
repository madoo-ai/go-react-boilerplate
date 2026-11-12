package xmodule

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"go-react-boilerplate/internal/adapter/mailer"
	pgadapter "go-react-boilerplate/internal/adapter/postgres"
	"go-react-boilerplate/internal/adapter/ratelimit"
	"go-react-boilerplate/internal/adapter/security"
	auth_usecase "go-react-boilerplate/internal/app/usecase/auth"
	dashboard_usecase "go-react-boilerplate/internal/app/usecase/dashboard"
	health_usecase "go-react-boilerplate/internal/app/usecase/health"
	roles_usecase "go-react-boilerplate/internal/app/usecase/roles"
	settings_usecase "go-react-boilerplate/internal/app/usecase/settings"
	users_usecase "go-react-boilerplate/internal/app/usecase/users"
	"go-react-boilerplate/internal/handler/httpapi"
	"go-react-boilerplate/internal/handler/webui"
)

// Dependencies carries everything the server needs to run. This is the only
// package that knows concrete adapters exist.
type Dependencies struct {
	API   *httpapi.Server
	WebUI *webui.Server
}

// InitDependencies applies migrations, opens the database, and wires
// adapters → use cases → handlers.
func InitDependencies(ctx context.Context, cfg Config, log *slog.Logger) (Dependencies, error) {
	if err := cfg.Validate(); err != nil {
		return Dependencies{}, err
	}
	if cfg.AuthSecret == "" {
		return Dependencies{}, fmt.Errorf("auth secret resolved empty")
	}
	if cfg.CookieSecure && cfg.AppURL == "http://localhost:5173" {
		return Dependencies{}, fmt.Errorf("COOKIE_SECURE=true requires an https APP_URL outside local development")
	}

	if err := pgadapter.Migrate(cfg.DSN); err != nil {
		return Dependencies{}, fmt.Errorf("database migrate: %w", err)
	}
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		TranslateError: true,
		Logger:         gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return Dependencies{}, fmt.Errorf("database open: %w", err)
	}
	log.Info("database ready")

	users := pgadapter.NewUserStore(db)
	roles := pgadapter.NewRoleStore(db)
	settings := pgadapter.NewSettingStore(db)
	sessions := pgadapter.NewSessionStore(db)
	verifications := pgadapter.NewVerificationStore(db)
	twoFactors := pgadapter.NewTwoFactorStore(db)
	dashboardStore := pgadapter.NewDashboardStore(db)
	healthStore := pgadapter.NewHealth(db)

	hasher := security.NewHasher()
	tokens := security.NewTokens()
	totpService := security.NewTOTP()
	cipher, err := security.NewCipher([]byte(cfg.AuthSecret))
	if err != nil {
		return Dependencies{}, fmt.Errorf("secret cipher: %w", err)
	}
	limiter := ratelimit.New(cfg.RateLimitPerMinute, time.Minute)
	mail := mailer.New(mailer.Config{
		Host: cfg.SMTPHost, Port: cfg.SMTPPort, User: cfg.SMTPUser,
		Password: cfg.SMTPPassword, From: cfg.SMTPFrom,
		AppName: cfg.AppName, Secure: cfg.SMTPSecure,
	}, log)

	clock := systemClock{}
	authUC := auth_usecase.New(
		users, roles, settings, sessions, verifications, twoFactors,
		hasher, tokens, totpService, cipher, mail, limiter, clock,
		auth_usecase.Config{AppURL: cfg.AppURL, SessionTTL: cfg.SessionTTL, Secret: []byte(cfg.AuthSecret)},
	)
	usersUC := users_usecase.New(users, roles, settings, sessions, hasher, clock)
	rolesUC := roles_usecase.New(roles, clock)
	settingsUC := settings_usecase.New(
		users, sessions, twoFactors, hasher, tokens, totpService, cipher,
		authUC, clock, cfg.AppName,
	)
	dashboardUC := dashboard_usecase.New(dashboardStore, clock)
	healthUC := health_usecase.New(healthStore, Version)

	// Expired sessions have no value; clear them at boot.
	if err := sessions.DeleteExpired(ctx, clock.Now().UTC()); err != nil {
		log.Warn("could not prune expired sessions", "err", err)
	}

	api := httpapi.New(authUC, usersUC, rolesUC, settingsUC, dashboardUC, healthUC, httpapi.Config{
		AppURL:        cfg.AppURL,
		SessionTTL:    cfg.SessionTTL,
		CookieSecure:  cfg.CookieSecure,
		TrustedOrigin: cfg.TrustedOrigin,
	}, log)

	return Dependencies{API: api, WebUI: webui.New(cfg.AppName, cfg.AppLocale)}, nil
}

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }
