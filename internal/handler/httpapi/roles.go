package httpapi

import (
	"net/http"

	"go-react-boilerplate/internal/app/usecase/roles"
)

type roleInputRequest struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

func (s *Server) handleListRoles(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	result, err := s.roles.List(r.Context(), session.Actor())
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	roleRows := make([]roleDTO, 0, len(result.Roles))
	for index := range result.Roles {
		roleRows = append(roleRows, *toRoleDTO(&result.Roles[index]))
	}
	permissionRows := make([]permissionDTO, 0, len(result.Permissions))
	for _, permission := range result.Permissions {
		permissionRows = append(permissionRows, permissionDTO{Key: string(permission.Key), Label: permission.Label})
	}
	writeJSON(w, http.StatusOK, roleListDTO{Roles: roleRows, Permissions: permissionRows})
}

func (s *Server) handleCreateRole(w http.ResponseWriter, r *http.Request) {
	var request roleInputRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	session := sessionFrom(r)
	input := roles.Input{
		Name: request.Name, Label: request.Label,
		Description: request.Description, Permissions: request.Permissions,
	}
	if err := s.roles.Create(r.Context(), session.Actor(), input); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]bool{"ok": true})
}

func (s *Server) handleUpdateRole(w http.ResponseWriter, r *http.Request) {
	var request roleInputRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	session := sessionFrom(r)
	input := roles.Input{
		ID: r.PathValue("id"), Name: request.Name, Label: request.Label,
		Description: request.Description, Permissions: request.Permissions,
	}
	if err := s.roles.Update(r.Context(), session.Actor(), input); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleDeleteRole(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	if err := s.roles.Delete(r.Context(), session.Actor(), r.PathValue("id")); err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
