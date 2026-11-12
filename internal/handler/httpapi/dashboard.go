package httpapi

import (
	"net/http"
)

type roleDatumDTO struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Users int    `json:"users"`
}

type signupDatumDTO struct {
	Date  string `json:"date"`
	Users int    `json:"users"`
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	session := sessionFrom(r)
	data, err := s.dashboard.Get(r.Context(), session)
	if err != nil {
		writeUseCaseError(w, s, err)
		return
	}
	roles := make([]roleDatumDTO, 0, len(data.RoleData))
	for _, role := range data.RoleData {
		roles = append(roles, roleDatumDTO{Name: role.Name, Label: role.Label, Users: role.Users})
	}
	signups := make([]signupDatumDTO, 0, len(data.SignupData))
	for _, signup := range data.SignupData {
		signups = append(signups, signupDatumDTO{Date: signup.Date, Users: signup.Users})
	}
	recent := make([]userDTO, 0, len(data.RecentUsers))
	for _, user := range data.RecentUsers {
		recent = append(recent, toUserDTO(user))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"totalUsers": data.TotalUsers, "activeUsers": data.ActiveUsers,
		"pendingUsers": data.PendingUsers, "rolesCount": data.RolesCount,
		"roleData": roles, "signupData": signups, "recentUsers": recent,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.health.Check(r.Context()))
}
