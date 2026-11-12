package httpapi

import (
	"time"

	"go-react-boilerplate/internal/entity"
)

type roleDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Description *string  `json:"description"`
	IsSystem    bool     `json:"isSystem"`
	Permissions []string `json:"permissions"`
	UserCount   int      `json:"userCount"`
}

type userDTO struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	EmailVerified    bool      `json:"emailVerified"`
	Image            *string   `json:"image"`
	Status           string    `json:"status"`
	PhoneNumber      *string   `json:"phoneNumber"`
	JobTitle         *string   `json:"jobTitle"`
	Department       *string   `json:"department"`
	TwoFactorEnabled bool      `json:"twoFactorEnabled"`
	Role             *roleDTO  `json:"role"`
	Permissions      []string  `json:"permissions"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type sessionDTO struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	IPAddress *string   `json:"ipAddress"`
	UserAgent *string   `json:"userAgent"`
}

type sessionResponse struct {
	Session sessionDTO `json:"session"`
	User    userDTO    `json:"user"`
}

type listUserDTO struct {
	userDTO
	CanManage bool `json:"canManage"`
	CanDelete bool `json:"canDelete"`
}

type roleListDTO struct {
	Roles       []roleDTO       `json:"roles"`
	Permissions []permissionDTO `json:"permissions"`
}

type permissionDTO struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type errorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func toNullable(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func toRoleDTO(role *entity.Role) *roleDTO {
	if role == nil {
		return nil
	}
	permissions := make([]string, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
		permissions = append(permissions, string(permission))
	}
	return &roleDTO{
		ID: role.ID, Name: role.Name, Label: role.Label,
		Description: toNullable(role.Description), IsSystem: role.IsSystem,
		Permissions: permissions, UserCount: role.UserCount,
	}
}

func toUserDTO(user entity.User) userDTO {
	permissions := make([]string, 0, len(user.Permissions))
	for _, permission := range user.Permissions {
		permissions = append(permissions, string(permission))
	}
	return userDTO{
		ID: user.ID, Name: user.Name, Email: user.Email, EmailVerified: user.EmailVerified,
		Image: toNullable(user.Image), Status: string(user.Status),
		PhoneNumber: toNullable(user.PhoneNumber), JobTitle: toNullable(user.JobTitle),
		Department: toNullable(user.Department), TwoFactorEnabled: user.TwoFactorEnabled,
		Role: toRoleDTO(user.Role), Permissions: permissions,
		CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
	}
}

func toSessionResponse(session *entity.SessionContext) sessionResponse {
	return sessionResponse{
		Session: sessionDTO{
			ID: session.Session.ID, CreatedAt: session.Session.CreatedAt,
			UpdatedAt: session.Session.UpdatedAt, ExpiresAt: session.Session.ExpiresAt,
			IPAddress: toNullable(session.Session.IPAddress),
			UserAgent: toNullable(session.Session.UserAgent),
		},
		User: toUserDTO(session.User),
	}
}
