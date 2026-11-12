// Package httpapi is the cookie-authenticated JSON driving adapter for the
// React app. It calls the use cases; it contains no business rules.
package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go-react-boilerplate/internal/app/usecase/auth"
	"go-react-boilerplate/internal/app/usecase/dashboard"
	"go-react-boilerplate/internal/app/usecase/health"
	"go-react-boilerplate/internal/app/usecase/roles"
	"go-react-boilerplate/internal/app/usecase/settings"
	"go-react-boilerplate/internal/app/usecase/users"
	"go-react-boilerplate/internal/entity"
)

const (
	sessionCookieName         = "gr_session"
	challengeCookieName       = "gr_2fa"
	passwordConfirmCookieName = "gr_password_confirmed" //nolint:gosec // This is a cookie name, not a credential.
	uiRequestHeader           = "X-UI-Request"
)

type Config struct {
	AppURL        string
	SessionTTL    time.Duration
	CookieSecure  bool
	TrustedOrigin []string
}

type Server struct {
	auth      auth.UseCase
	users     users.UseCase
	roles     roles.UseCase
	settings  settings.UseCase
	dashboard dashboard.UseCase
	health    health.UseCase
	log       *slog.Logger
	cfg       Config
}

func New(
	authUC auth.UseCase,
	usersUC users.UseCase,
	rolesUC roles.UseCase,
	settingsUC settings.UseCase,
	dashboardUC dashboard.UseCase,
	healthUC health.UseCase,
	cfg Config,
	log *slog.Logger,
) *Server {
	if cfg.SessionTTL <= 0 {
		cfg.SessionTTL = 7 * 24 * time.Hour
	}
	return &Server{
		auth: authUC, users: usersUC, roles: rolesUC, settings: settingsUC,
		dashboard: dashboardUC, health: healthUC, cfg: cfg, log: log,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("POST /api/auth/register", s.handleRegister)
	mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
	mux.HandleFunc("POST /api/auth/verify-email", s.handleVerifyEmail)
	mux.HandleFunc("POST /api/auth/forgot-password", s.handleForgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", s.handleResetPassword)
	mux.HandleFunc("POST /api/auth/two-factor/verify", s.handleTwoFactorVerify)

	mux.HandleFunc("GET /api/auth/session", s.requireSession(s.handleSession))
	mux.HandleFunc("POST /api/auth/verify-email/request", s.requireSession(s.handleRequestVerification))
	mux.HandleFunc("POST /api/auth/confirm-password", s.requireSession(s.handleConfirmPassword))

	mux.HandleFunc("GET /api/dashboard", s.requirePermission(entity.PermissionViewAdmin, s.handleDashboard))

	mux.HandleFunc("GET /api/users", s.requirePermission(entity.PermissionManageUsers, s.handleListUsers))
	mux.HandleFunc("POST /api/users", s.requirePermission(entity.PermissionManageUsers, s.handleCreateUser))
	mux.HandleFunc("GET /api/users/{id}", s.requirePermission(entity.PermissionManageUsers, s.handleGetUser))
	mux.HandleFunc("PATCH /api/users/{id}", s.requirePermission(entity.PermissionManageUsers, s.handleUpdateUser))
	mux.HandleFunc("DELETE /api/users/{id}", s.requirePermission(entity.PermissionManageUsers, s.handleDeleteUser))
	mux.HandleFunc("POST /api/users/{id}/verify", s.requirePermission(entity.PermissionManageUsers, s.handleVerifyUser))
	mux.HandleFunc("PATCH /api/admin/registration-approval", s.requirePermission(entity.PermissionManageUsers, s.handleRegistrationApproval))

	mux.HandleFunc("GET /api/roles", s.requirePermission(entity.PermissionManageRoles, s.handleListRoles))
	mux.HandleFunc("POST /api/roles", s.requirePermission(entity.PermissionManageRoles, s.handleCreateRole))
	mux.HandleFunc("PATCH /api/roles/{id}", s.requirePermission(entity.PermissionManageRoles, s.handleUpdateRole))
	mux.HandleFunc("DELETE /api/roles/{id}", s.requirePermission(entity.PermissionManageRoles, s.handleDeleteRole))

	mux.HandleFunc("GET /api/settings/account", s.requireSession(s.handleSettingsAccount))
	mux.HandleFunc("PATCH /api/settings/profile", s.requireSession(s.handleUpdateProfile))
	mux.HandleFunc("POST /api/settings/password", s.requireSession(s.handleChangePassword))
	mux.HandleFunc("POST /api/settings/two-factor/enable", s.requireSession(s.handleEnableTwoFactor))
	mux.HandleFunc("POST /api/settings/two-factor/confirm", s.requireSession(s.handleConfirmTwoFactor))
	mux.HandleFunc("POST /api/settings/two-factor/disable", s.requireSession(s.handleDisableTwoFactor))
	mux.HandleFunc("DELETE /api/settings/account", s.requireSession(s.handleDeleteAccount))

	mux.HandleFunc("GET /api/", func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, entity.NewError("errors.notFound"))
	})

	return s.securityHeaders(s.mutationGuard(mux))
}

// --- middleware -----------------------------------------------------------

type ctxKey string

const sessionKey ctxKey = "session"

func (s *Server) requireSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.currentSession(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, entity.ErrUnauthorized)
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), sessionKey, session)))
	}
}

func (s *Server) requirePermission(permission entity.Permission, next http.HandlerFunc) http.HandlerFunc {
	return s.requireSession(func(w http.ResponseWriter, r *http.Request) {
		session := sessionFrom(r)
		if session == nil || !session.Actor().Can(permission) {
			writeError(w, http.StatusForbidden, entity.ErrForbidden)
			return
		}
		next(w, r)
	})
}

func (s *Server) currentSession(r *http.Request) (*entity.SessionContext, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, entity.ErrUnauthorized
	}
	return s.auth.Session(r.Context(), cookie.Value)
}

func sessionFrom(r *http.Request) *entity.SessionContext {
	session, _ := r.Context().Value(sessionKey).(*entity.SessionContext)
	return session
}

// mutationGuard protects state-changing requests. A custom header cannot be
// set by a cross-origin form post, and the Origin check rejects a caller that
// does manage a preflight.
func (s *Server) mutationGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get(uiRequestHeader) != "1" {
			writeError(w, http.StatusForbidden, entity.ErrForbidden)
			return
		}
		if origin := r.Header.Get("Origin"); origin != "" && !s.originAllowed(origin, r.Host) {
			writeError(w, http.StatusForbidden, entity.ErrForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) originAllowed(origin, host string) bool {
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if strings.EqualFold(parsed.Host, host) {
		return true
	}
	if s.cfg.AppURL != "" {
		if app, err := url.Parse(s.cfg.AppURL); err == nil && strings.EqualFold(parsed.Host, app.Host) {
			return true
		}
	}
	for _, allowed := range s.cfg.TrustedOrigin {
		if strings.EqualFold(origin, strings.TrimRight(allowed, "/")) {
			return true
		}
	}
	return false
}

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

// --- helpers --------------------------------------------------------------

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		writeError(w, http.StatusBadRequest, entity.NewError("errors.invalidRequest"))
		return false
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, entity.NewError("errors.invalidRequest"))
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	code := "errors.unexpected"
	fields := map[string]string(nil)
	var coded *entity.CodedError
	if errors.As(err, &coded) {
		code = coded.Code
		fields = coded.Fields
	}
	writeJSON(w, status, errorResponse{Error: code, Fields: fields})
}

func statusFor(err error) int {
	var coded *entity.CodedError
	if !errors.As(err, &coded) {
		return http.StatusInternalServerError
	}
	switch coded.Code {
	case "errors.unauthorized", "auth.messages.invalidCredentials":
		return http.StatusUnauthorized
	case "errors.forbidden", "errors.confirmPasswordRequired",
		"errors.cannotDeleteSelf", "errors.lastSuperadmin", "errors.protectedUser",
		"errors.immutableRole", "errors.defaultRole",
		"auth.messages.pendingApproval", "auth.messages.suspended":
		return http.StatusForbidden
	case "errors.userNotFound", "errors.roleNotFound", "errors.notFound":
		return http.StatusNotFound
	case "errors.emailInUse", "errors.roleInUse":
		return http.StatusConflict
	case "errors.tooManyRequests":
		return http.StatusTooManyRequests
	case "errors.unexpected":
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

func writeUseCaseError(w http.ResponseWriter, s *Server, err error) {
	status := statusFor(err)
	if status == http.StatusInternalServerError {
		s.log.Error("request failed", "err", err)
	}
	writeError(w, status, err)
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (s *Server) setSessionCookie(w http.ResponseWriter, token string, remember bool) {
	//nolint:gosec // HttpOnly, SameSite and Secure are set; Secure follows APP_URL.
	cookie := &http.Cookie{
		Name: sessionCookieName, Value: token, Path: "/",
		HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: s.cfg.CookieSecure,
	}
	if remember {
		cookie.MaxAge = int(s.cfg.SessionTTL.Seconds())
	}
	http.SetCookie(w, cookie)
}

func (s *Server) clearSessionCookie(w http.ResponseWriter) {
	//nolint:gosec // HttpOnly, SameSite and Secure are set; Secure follows APP_URL.
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookieName, Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: s.cfg.CookieSecure,
	})
}

func (s *Server) setChallengeCookie(w http.ResponseWriter, value string) {
	//nolint:gosec // HttpOnly, SameSite and Secure are set; Secure follows APP_URL.
	http.SetCookie(w, &http.Cookie{
		Name: challengeCookieName, Value: value, Path: "/", MaxAge: 5 * 60,
		HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: s.cfg.CookieSecure,
	})
}

func (s *Server) clearChallengeCookie(w http.ResponseWriter) {
	//nolint:gosec // HttpOnly, SameSite and Secure are set; Secure follows APP_URL.
	http.SetCookie(w, &http.Cookie{
		Name: challengeCookieName, Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: s.cfg.CookieSecure,
	})
}

func (s *Server) setPasswordConfirmCookie(w http.ResponseWriter, value string) {
	//nolint:gosec // HttpOnly, SameSite and Secure are set; Secure follows APP_URL.
	http.SetCookie(w, &http.Cookie{
		Name: passwordConfirmCookieName, Value: value, Path: "/", MaxAge: 15 * 60,
		HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: s.cfg.CookieSecure,
	})
}

func (s *Server) clearPasswordConfirmCookie(w http.ResponseWriter) {
	//nolint:gosec // HttpOnly, SameSite and Secure are set; Secure follows APP_URL.
	http.SetCookie(w, &http.Cookie{
		Name: passwordConfirmCookieName, Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: s.cfg.CookieSecure,
	})
}
