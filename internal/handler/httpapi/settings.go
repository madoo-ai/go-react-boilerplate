package httpapi

import (
	"net/http"

	"go-react-boilerplate/internal/app/usecase/settings"
)

type profileRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	JobTitle   string `json:"jobTitle"`
	Department string `json:"department"`
}

func (s *Server) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	var request profileRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	err := s.settings.UpdateProfile(r.Context(), session, settings.ProfileInput{
		Name: request.Name, Email: request.Email, Phone: request.Phone,
		JobTitle: request.JobTitle, Department: request.Department,
	})
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	updated, err := s.currentSession(r)
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, toSessionResponse(updated))
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (s *Server) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	var request changePasswordRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	err := s.settings.ChangePassword(r.Context(), session, settings.PasswordInput{
		Current: request.CurrentPassword, Next: request.NewPassword,
	})
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleSettingsAccount(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	overview, err := s.settings.AccountOverview(r.Context(), session, s.hasPasswordConfirmation(r))
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{
		"canDelete": overview.CanDelete, "confirmed": overview.Confirmed,
	})
}

type passwordRequest struct {
	Password string `json:"password"`
}

func (s *Server) handleEnableTwoFactor(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	var request passwordRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	setup, err := s.settings.EnableTwoFactor(r.Context(), session, request.Password)
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"secret": setup.Secret, "uri": setup.URI, "qrDataUrl": setup.QRDataURL,
		"backupCodes": setup.BackupCodes,
	})
}

type twoFactorConfirmRequest struct {
	Code string `json:"code"`
}

func (s *Server) handleConfirmTwoFactor(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	var request twoFactorConfirmRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := s.settings.ConfirmTwoFactor(r.Context(), session, request.Code); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleDisableTwoFactor(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	var request passwordRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := s.settings.DisableTwoFactor(r.Context(), session, request.Password); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	if err := s.settings.DeleteAccount(r.Context(), session, s.hasPasswordConfirmation(r)); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	s.clearSessionCookie(w)
	s.clearChallengeCookie(w)
	s.clearPasswordConfirmCookie(w)
	w.WriteHeader(http.StatusNoContent)
}
