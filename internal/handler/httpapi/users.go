package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"go-react-boilerplate/internal/app/usecase/users"
	"go-react-boilerplate/internal/entity"
)

type userInputRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	JobTitle   string `json:"jobTitle"`
	Department string `json:"department"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	Password   string `json:"password"`
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	if session == nil {
		writeError(w, http.StatusUnauthorized, entity.ErrUnauthorized)
		return
	}
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	filters := entity.UserFilters{
		Page: page, Search: strings.TrimSpace(query.Get("search")),
		Role: strings.TrimSpace(query.Get("role")), Status: strings.TrimSpace(query.Get("status")),
	}
	result, err := s.users.List(r.Context(), session.Actor(), filters)
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	rows := make([]listUserDTO, 0, len(result.Users))
	for _, item := range result.Users {
		rows = append(rows, listUserDTO{
			userDTO: toUserDTO(item.User), CanManage: item.CanManage, CanDelete: item.CanDelete,
		})
	}
	roles := make([]roleDTO, 0, len(result.Roles))
	for index := range result.Roles {
		roles = append(roles, *toRoleDTO(&result.Roles[index]))
	}
	assignable := make([]roleDTO, 0, len(result.AssignableRoles))
	for index := range result.AssignableRoles {
		assignable = append(assignable, *toRoleDTO(&result.AssignableRoles[index]))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"users": rows, "total": result.Total, "page": result.Page,
		"pageCount": result.PageCount, "pageSize": result.PageSize,
		"roles": roles, "assignableRoles": assignable,
		"registrationApproval": result.RegistrationApproval,
		"superadminCount":      result.SuperadminCount,
	})
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var request userInputRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	session := sessionFrom(r)
	input := users.Input{
		Name: request.Name, Email: request.Email, Phone: request.Phone,
		JobTitle: request.JobTitle, Department: request.Department,
		Role: request.Role, Status: request.Status, Password: request.Password,
	}
	if err := s.users.Create(r.Context(), session.Actor(), input); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]bool{"ok": true})
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	result, err := s.users.GetForEdit(r.Context(), session.Actor(), r.PathValue("id"))
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	assignable := make([]roleDTO, 0, len(result.AssignableRoles))
	for index := range result.AssignableRoles {
		assignable = append(assignable, *toRoleDTO(&result.AssignableRoles[index]))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user": toUserDTO(result.User), "assignableRoles": assignable, "canManage": result.CanManage,
	})
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var request userInputRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	session := sessionFrom(r)
	input := users.Input{
		ID: r.PathValue("id"), Name: request.Name, Email: request.Email,
		Phone: request.Phone, JobTitle: request.JobTitle, Department: request.Department,
		Role: request.Role, Status: request.Status, Password: request.Password,
	}
	if err := s.users.Update(r.Context(), session.Actor(), input); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	if err := s.users.Delete(r.Context(), session.Actor(), r.PathValue("id")); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleVerifyUser(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	if err := s.users.Verify(r.Context(), session.Actor(), r.PathValue("id")); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

type registrationApprovalRequest struct {
	Enabled bool `json:"enabled"`
}

func (s *Server) handleRegistrationApproval(w http.ResponseWriter, r *http.Request) {
	var request registrationApprovalRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	session := sessionFrom(r)
	if err := s.users.SetRegistrationApproval(r.Context(), session.Actor(), request.Enabled); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
