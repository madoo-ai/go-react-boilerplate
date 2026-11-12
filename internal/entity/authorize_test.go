package entity

import "testing"

func actor(role string, permissions ...Permission) Actor {
	return Actor{ID: "actor-1", RoleName: role, Permissions: permissions}
}

func TestCanManageUser(t *testing.T) {
	hr := actor(RoleHR, PermissionManageUsers)
	superadmin := actor(RoleSuperadmin, PermissionManageUsers, PermissionManageRoles, PermissionViewAdmin)
	regular := actor(RoleUser)

	regularTarget := &User{ID: "u-1", Role: &Role{Name: RoleUser}}
	superTarget := &User{ID: "u-2", Role: &Role{Name: RoleSuperadmin}}

	cases := []struct {
		name   string
		actor  Actor
		target *User
		want   bool
	}{
		{"HR can manage a regular user", hr, regularTarget, true},
		{"HR cannot manage a superadmin", hr, superTarget, false},
		{"superadmin can manage a superadmin", superadmin, superTarget, true},
		{"regular user cannot manage anyone", regular, regularTarget, false},
	}
	for _, test := range cases {
		if got := CanManageUser(test.actor, test.target); got != test.want {
			t.Errorf("%s: got %v, want %v", test.name, got, test.want)
		}
	}
}

func TestCanDeleteUser(t *testing.T) {
	superadmin := actor(RoleSuperadmin, PermissionManageUsers, PermissionManageRoles)
	target := &User{ID: "u-2", Role: &Role{Name: RoleSuperadmin}}
	self := &User{ID: superadmin.ID, Role: &Role{Name: RoleSuperadmin}}

	if CanDeleteUser(superadmin, self, 2) {
		t.Fatal("an actor must not delete itself")
	}
	if CanDeleteUser(superadmin, target, 1) {
		t.Fatal("the last superadmin must not be deleted")
	}
	if !CanDeleteUser(superadmin, target, 2) {
		t.Fatal("a superadmin should be deletable while another remains")
	}
}

func TestCanAssignRole(t *testing.T) {
	hr := actor(RoleHR, PermissionManageUsers)
	superadmin := actor(RoleSuperadmin, PermissionManageUsers, PermissionManageRoles)

	if hr.Can(PermissionManageRoles) {
		t.Fatal("HR must not manage roles")
	}
	if CanAssignRole(hr, RoleSuperadmin) {
		t.Fatal("HR must not assign superadmin")
	}
	if !CanAssignRole(hr, RoleUser) {
		t.Fatal("HR should assign the default role")
	}
	if !CanAssignRole(superadmin, RoleSuperadmin) {
		t.Fatal("a superadmin should assign superadmin")
	}
}

func TestCanManageRole(t *testing.T) {
	superadmin := actor(RoleSuperadmin, PermissionManageRoles)
	hr := actor(RoleHR, PermissionManageUsers)

	if CanManageRole(hr, RoleUser) {
		t.Fatal("HR must not manage roles")
	}
	if CanManageRole(superadmin, RoleSuperadmin) {
		t.Fatal("the superadmin role is immutable")
	}
	if !CanManageRole(superadmin, RoleHR) {
		t.Fatal("a superadmin should manage custom roles")
	}
}

func TestValidationHelpers(t *testing.T) {
	if !ValidPassword("password") || ValidPassword("short") {
		t.Fatal("password bounds are wrong")
	}
	if !ValidEmail("person@example.com") || ValidEmail("not-an-email") {
		t.Fatal("email validation is wrong")
	}
	if !ValidRoleName("support_lead") || ValidRoleName("Support Lead") {
		t.Fatal("role name validation is wrong")
	}
	if NormalizeEmail("  Person@Example.COM ") != "person@example.com" {
		t.Fatal("email normalization is wrong")
	}
}
