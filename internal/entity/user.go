package entity

import (
	"regexp"
	"strings"
	"time"
)

// UserStatus is the account lifecycle state.
type UserStatus string

const (
	StatusPending   UserStatus = "pending"
	StatusActive    UserStatus = "active"
	StatusSuspended UserStatus = "suspended"
)

func IsUserStatus(value string) bool {
	switch UserStatus(value) {
	case StatusPending, StatusActive, StatusSuspended:
		return true
	default:
		return false
	}
}

// User is one account. Role is nil when no role has been assigned.
type User struct {
	ID               string
	Name             string
	Email            string
	EmailVerified    bool
	Image            string
	RoleID           string
	Status           UserStatus
	PhoneNumber      string
	JobTitle         string
	Department       string
	TwoFactorEnabled bool
	Role             *Role
	Permissions      []Permission
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Actor is the authenticated principal used by RBAC policies. It is built
// from the session and never trusts client input.
type Actor struct {
	ID          string
	Name        string
	Email       string
	RoleName    string
	Permissions []Permission
}

func (a Actor) Can(permission Permission) bool {
	for _, granted := range a.Permissions {
		if granted == permission {
			return true
		}
	}
	return false
}

func (a Actor) IsSuperadmin() bool { return a.RoleName == RoleSuperadmin }

// UserFilters is the validated query for the users list.
type UserFilters struct {
	Page   int
	Search string
	Role   string
	Status string
}

const UsersPageSize = 10

// NormalizeEmail lower-cases and trims an address. All writes go through it
// so uniqueness is case-insensitive in practice.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func ValidEmail(email string) bool { return emailPattern.MatchString(email) }

// ValidateName enforces the same 2..120 rule on both client and server.
func ValidateName(name string) bool {
	trimmed := strings.TrimSpace(name)
	return len([]rune(trimmed)) >= 2 && len([]rune(trimmed)) <= 120
}

func ValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 128
}

var roleNamePattern = regexp.MustCompile(`^[a-z0-9_]+$`)

func ValidRoleName(name string) bool {
	return len(name) >= 2 && len(name) <= 60 && roleNamePattern.MatchString(name)
}
