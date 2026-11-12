package entity

// RBAC policies. These are pure functions: no I/O, no HTTP, no storage. The
// usecases call them close to the data so a client cannot bypass a rule that
// only exists in the UI.

func Can(actor Actor, permission Permission) bool { return actor.Can(permission) }

// CanManageUser lets HR manage regular users but keeps superadmins out of
// reach unless the actor is a superadmin.
func CanManageUser(actor Actor, target *User) bool {
	if !actor.Can(PermissionManageUsers) {
		return false
	}
	if target != nil && target.Role != nil && target.Role.Name == RoleSuperadmin && !actor.IsSuperadmin() {
		return false
	}
	return true
}

// CanDeleteUser adds self-deletion and last-superadmin protection.
func CanDeleteUser(actor Actor, target *User, superadminCount int) bool {
	if !CanManageUser(actor, target) {
		return false
	}
	if actor.ID == target.ID {
		return false
	}
	if target.Role != nil && target.Role.Name == RoleSuperadmin && superadminCount <= 1 {
		return false
	}
	return true
}

// CanAssignRole prevents privilege escalation: only a manager of roles may
// hand out the superadmin role.
func CanAssignRole(actor Actor, roleName string) bool {
	if roleName == RoleSuperadmin {
		return actor.Can(PermissionManageRoles)
	}
	return actor.Can(PermissionManageUsers)
}

// CanManageRole protects the immutable superadmin role. The default user role
// is editable but cannot be renamed or deleted (checked in the usecase).
func CanManageRole(actor Actor, roleName string) bool {
	if !actor.Can(PermissionManageRoles) {
		return false
	}
	return roleName != RoleSuperadmin
}
