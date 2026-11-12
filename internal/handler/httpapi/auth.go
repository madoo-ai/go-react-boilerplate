package httpapi

import (
	"net/http"
	"time"

	"go-react-boilerplate/internal/app/usecase/auth"
	"go-react-boilerplate/internal/entity"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var request registerRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := s.auth.Register(r.Context(), auth.RegisterInput{
		Name: request.Name, Email: request.Email, Password: request.Password, IP: clientIP(r),
	})
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"status": result.Status, "emailSent": result.EmailSent,
	})
}

type loginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := s.auth.Login(r.Context(), auth.LoginInput{
		Email: request.Email, Password: request.Password, IP: clientIP(r),
		UserAgent: r.UserAgent(), RememberMe: request.RememberMe,
	})
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	if result.TwoFactorRequired {
		s.setChallengeCookie(w, result.Challenge)
		writeJSON(w, http.StatusOK, map[string]any{"twoFactorRequired": true})
		return
	}
	s.setSessionCookie(w, result.Token, request.RememberMe)
	writeJSON(w, http.StatusOK, toSessionResponse(result.Context))
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		_ = s.auth.Logout(r.Context(), cookie.Value)
	}
	s.clearSessionCookie(w)
	s.clearChallengeCookie(w)
	s.clearPasswordConfirmCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	if session == nil {
		writeError(w, http.StatusUnauthorized, entity.ErrUnauthorized)
		return
	}
	writeJSON(w, http.StatusOK, toSessionResponse(session))
}

type tokenRequest struct {
	Token string `json:"token"`
}

func (s *Server) handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	var request tokenRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := s.auth.VerifyEmail(r.Context(), request.Token); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleRequestVerification(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	if session == nil {
		writeError(w, http.StatusUnauthorized, entity.ErrUnauthorized)
		return
	}
	if err := s.auth.RequestEmailVerification(r.Context(), session.User.ID, clientIP(r)); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (s *Server) handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var request forgotPasswordRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := s.auth.ForgotPassword(r.Context(), request.Email, clientIP(r)); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

type resetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (s *Server) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var request resetPasswordRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := s.auth.ResetPassword(r.Context(), request.Token, request.Password); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

type twoFactorRequest struct {
	Code   string `json:"code"`
	Backup bool   `json:"backup"`
}

func (s *Server) handleTwoFactorVerify(w http.ResponseWriter, r *http.Request) {
	challenge, err := r.Cookie(challengeCookieName)
	if err != nil || challenge.Value == "" {
		writeError(w, http.StatusUnauthorized, entity.ErrUnauthorized)
		return
	}
	var request twoFactorRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := s.auth.VerifyTwoFactor(r.Context(), challenge.Value, request.Code, request.Backup, clientIP(r), r.UserAgent())
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	s.clearChallengeCookie(w)
	s.setSessionCookie(w, result.Token, true)
	writeJSON(w, http.StatusOK, toSessionResponse(result.Context))
}

type confirmPasswordRequest struct {
	Password string `json:"password"`
}

func (s *Server) handleConfirmPassword(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	if session == nil {
		writeError(w, http.StatusUnauthorized, entity.ErrUnauthorized)
		return
	}
	var request confirmPasswordRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	value, err := s.auth.ConfirmPassword(r.Context(), session, request.Password)
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	s.setPasswordConfirmCookie(w, value)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) hasPasswordConfirmation(r *http.Request) bool {
	session := sessionFrom(r)
	if session == nil {
		return false
	}
	cookie, err := r.Cookie(passwordConfirmCookieName)
	if err != nil {
		return false
	}
	return s.auth.HasRecentPasswordConfirmation(session.User.ID, cookie.Value, time.Now().UTC())
}
