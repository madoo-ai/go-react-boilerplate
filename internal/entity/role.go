package entity

import (
	"time"
)

// Permission is a coarse capability granted to a role.
type Permission string

const (
	PermissionViewAdmin   Permission = "view_admin"
	PermissionManageUsers Permission = "manage_users"
	PermissionManageRoles Permission = "manage_roles"
)

// AllPermissions is the registry, in display order.
var AllPermissions = []Permission{
	PermissionViewAdmin,
	PermissionManageUsers,
	PermissionManageRoles,
}

var PermissionLabels = map[Permission]string{
	PermissionViewAdmin:   "Access the admin area",
	PermissionManageUsers: "Create, edit, verify and delete users",
	PermissionManageRoles: "Create, edit and delete roles",
}

func IsPermission(value string) bool {
	for _, permission := range AllPermissions {
		if string(permission) == value {
			return true
		}
	}
	return false
}

// Role names that the seed guarantees. Custom roles may be created at runtime.
const (
	RoleSuperadmin = "superadmin"
	RoleHR         = "hr"
	RoleUser       = "user"
)

var SystemRoleNames = []string{RoleSuperadmin, RoleHR, RoleUser}

type RoleDefinition struct {
	Name        string
	Label       string
	Description string
	IsSystem    bool
	Permissions []Permission
}

// RoleDefinitions is the seeded role catalogue. It is the single source of
// truth used by the seed command and by tests.
var RoleDefinitions = []RoleDefinition{
	{
		Name: RoleSuperadmin, Label: "Super Admin",
		Description: "Full access, including role management.",
		IsSystem:    true,
		Permissions: []Permission{PermissionViewAdmin, PermissionManageUsers, PermissionManageRoles},
	},
	{
		Name: RoleHR, Label: "HR Staff",
		Description: "Can manage users but cannot manage roles.",
		Permissions: []Permission{PermissionViewAdmin, PermissionManageUsers},
	},
	{
		Name: RoleUser, Label: "User",
		Description: "Self-service access only.",
		Permissions: []Permission{},
	},
}

func IsSystemRole(name string) bool {
	return name == RoleSuperadmin || name == RoleHR || name == RoleUser
}

// Role is an assignable role with its permissions and member count.
type Role struct {
	ID          string
	Name        string
	Label       string
	Description string
	IsSystem    bool
	Permissions []Permission
	UserCount   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
