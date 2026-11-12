package postgres

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"go-react-boilerplate/internal/entity"
)

// Storage models. Gorm types never leave this package; everything crosses as
// entities.

type roleRow struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name"`
	Label       string    `gorm:"column:label"`
	Description *string   `gorm:"column:description"`
	IsSystem    bool      `gorm:"column:is_system"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (roleRow) TableName() string { return "roles" }

type permissionRow struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	Key         string    `gorm:"column:key"`
	Label       string    `gorm:"column:label"`
	Description *string   `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (permissionRow) TableName() string { return "permissions" }

type rolePermissionRow struct {
	RoleID       string `gorm:"column:role_id;type:uuid;primaryKey"`
	PermissionID string `gorm:"column:permission_id;type:uuid;primaryKey"`
}

func (rolePermissionRow) TableName() string { return "role_permissions" }

type userRow struct {
	ID               string    `gorm:"column:id;type:uuid;primaryKey"`
	Name             string    `gorm:"column:name"`
	Email            string    `gorm:"column:email"`
	EmailVerified    bool      `gorm:"column:email_verified"`
	Image            *string   `gorm:"column:image"`
	RoleID           *string   `gorm:"column:role_id"`
	Status           string    `gorm:"column:status"`
	PhoneNumber      *string   `gorm:"column:phone_number"`
	JobTitle         *string   `gorm:"column:job_title"`
	Department       *string   `gorm:"column:department"`
	TwoFactorEnabled bool      `gorm:"column:two_factor_enabled"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (userRow) TableName() string { return "users" }

type localCredentialRow struct {
	UserID       string    `gorm:"column:user_id;type:uuid;primaryKey"`
	PasswordHash string    `gorm:"column:password_hash"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (localCredentialRow) TableName() string { return "local_credentials" }

type sessionRow struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey"`
	UserID    string    `gorm:"column:user_id;type:uuid"`
	TokenHash string    `gorm:"column:token_hash"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	IPAddress *string   `gorm:"column:ip_address"`
	UserAgent *string   `gorm:"column:user_agent"`
}

func (sessionRow) TableName() string { return "sessions" }

type settingRow struct {
	Key       string    `gorm:"column:key;primaryKey"`
	Value     *string   `gorm:"column:value"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (settingRow) TableName() string { return "settings" }

type verificationRow struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey"`
	Identifier string    `gorm:"column:identifier"`
	Purpose    string    `gorm:"column:purpose"`
	TokenHash  string    `gorm:"column:token_hash"`
	ExpiresAt  time.Time `gorm:"column:expires_at"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (verificationRow) TableName() string { return "verifications" }

type twoFactorRow struct {
	ID                      string     `gorm:"column:id;type:uuid;primaryKey"`
	UserID                  string     `gorm:"column:user_id;type:uuid"`
	Secret                  string     `gorm:"column:secret"`
	BackupCodes             string     `gorm:"column:backup_codes"`
	Verified                bool       `gorm:"column:verified"`
	FailedVerificationCount int        `gorm:"column:failed_verification_count"`
	LockedUntil             *time.Time `gorm:"column:locked_until"`
	CreatedAt               time.Time  `gorm:"column:created_at"`
	UpdatedAt               time.Time  `gorm:"column:updated_at"`
}

func (twoFactorRow) TableName() string { return "two_factors" }

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func ptr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func toRole(row roleRow) entity.Role {
	return entity.Role{
		ID: row.ID, Name: row.Name, Label: row.Label,
		Description: deref(row.Description), IsSystem: row.IsSystem,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}

func toUser(row userRow, role *entity.Role, permissions []entity.Permission) entity.User {
	user := entity.User{
		ID: row.ID, Name: row.Name, Email: row.Email, EmailVerified: row.EmailVerified,
		Image: deref(row.Image), Status: entity.UserStatus(row.Status),
		PhoneNumber: deref(row.PhoneNumber), JobTitle: deref(row.JobTitle),
		Department: deref(row.Department), TwoFactorEnabled: row.TwoFactorEnabled,
		RoleID: deref(row.RoleID), Role: role, Permissions: permissions,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
	return user
}

func isNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }

func isUnique(err error) bool {
	return err != nil && (errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(err.Error(), "SQLSTATE 23505") ||
		strings.Contains(err.Error(), "duplicate key value"))
}
