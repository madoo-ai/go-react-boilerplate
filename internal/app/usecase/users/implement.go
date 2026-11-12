package users

import (
	"context"
	"errors"
	"strings"

	"go-react-boilerplate/internal/app/port"
	"go-react-boilerplate/internal/entity"
)

func (u *useCase) List(ctx context.Context, actor entity.Actor, filters entity.UserFilters) (ListResult, error) {
	if !actor.Can(entity.PermissionManageUsers) {
		return ListResult{}, entity.ErrForbidden
	}
	if filters.Page < 1 {
		filters.Page = 1
	}
	users, total, err := u.users.List(ctx, filters)
	if err != nil {
		return ListResult{}, err
	}
	pageCount := (total + entity.UsersPageSize - 1) / entity.UsersPageSize
	if pageCount < 1 {
		pageCount = 1
	}
	page := filters.Page
	if page > pageCount {
		page = pageCount
	}
	roles, err := u.roles.List(ctx)
	if err != nil {
		return ListResult{}, err
	}
	assignable := make([]entity.Role, 0, len(roles))
	for _, role := range roles {
		if entity.CanAssignRole(actor, role.Name) {
			assignable = append(assignable, role)
		}
	}
	approval, err := u.registrationApproval(ctx)
	if err != nil {
		return ListResult{}, err
	}
	superadmins, err := u.users.CountSuperadmins(ctx)
	if err != nil {
		return ListResult{}, err
	}

	items := make([]UserListItem, 0, len(users))
	for index := range users {
		user := users[index]
		items = append(items, UserListItem{
			User:      user,
			CanManage: entity.CanManageUser(actor, &user),
			CanDelete: entity.CanDeleteUser(actor, &user, superadmins),
		})
	}

	return ListResult{
		Users: items, Total: total, Page: page, PageCount: pageCount,
		PageSize: entity.UsersPageSize, Roles: roles, AssignableRoles: assignable,
		RegistrationApproval: approval, SuperadminCount: superadmins,
	}, nil
}

func (u *useCase) GetForEdit(ctx context.Context, actor entity.Actor, id string) (EditResult, error) {
	if !actor.Can(entity.PermissionManageUsers) {
		return EditResult{}, entity.ErrForbidden
	}
	target, err := u.users.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return EditResult{}, entity.ErrUserNotFound
		}
		return EditResult{}, err
	}
	if !entity.CanManageUser(actor, target) {
		return EditResult{}, entity.ErrForbidden
	}
	roles, err := u.roles.List(ctx)
	if err != nil {
		return EditResult{}, err
	}
	assignable := make([]entity.Role, 0, len(roles))
	for _, role := range roles {
		if entity.CanAssignRole(actor, role.Name) {
			assignable = append(assignable, role)
		}
	}
	return EditResult{User: *target, AssignableRoles: assignable, CanManage: true}, nil
}

func (u *useCase) Create(ctx context.Context, actor entity.Actor, input Input) error {
	if !actor.Can(entity.PermissionManageUsers) {
		return entity.ErrForbidden
	}
	fields, name, email := validate(input, true)
	if len(fields) > 0 {
		return entity.ValidationError(fields)
	}
	role, err := u.assignableRole(ctx, actor, input.Role)
	if err != nil {
		return err
	}
	if err := u.ensureEmailAvailable(ctx, email, ""); err != nil {
		return err
	}
	status := entity.StatusActive
	if entity.IsUserStatus(input.Status) {
		status = entity.UserStatus(input.Status)
	}
	hash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	now := u.clock.Now().UTC()
	user := &entity.User{
		ID: entity.NewUUID(), Name: name, Email: email, Status: status,
		RoleID: role.ID, PhoneNumber: clean(input.Phone), JobTitle: clean(input.JobTitle),
		Department: clean(input.Department), CreatedAt: now, UpdatedAt: now,
	}
	if err := u.users.CreateWithPassword(ctx, user, hash); err != nil {
		if errors.Is(err, port.ErrUniqueViolation) {
			return entity.ErrEmailInUse
		}
		return err
	}
	return nil
}

func (u *useCase) Update(ctx context.Context, actor entity.Actor, input Input) error {
	if !actor.Can(entity.PermissionManageUsers) {
		return entity.ErrForbidden
	}
	fields, name, email := validate(input, false)
	if len(fields) > 0 {
		return entity.ValidationError(fields)
	}
	target, err := u.users.GetByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrUserNotFound
		}
		return err
	}
	if !entity.CanManageUser(actor, target) {
		return entity.ErrForbidden
	}
	role, err := u.assignableRole(ctx, actor, input.Role)
	if err != nil {
		return err
	}
	status := target.Status
	if entity.IsUserStatus(input.Status) {
		status = entity.UserStatus(input.Status)
	}
	if target.ID == actor.ID && (role.Name != roleName(target.Role) || status != target.Status) {
		return entity.ErrProtectedUser
	}
	if target.Role != nil && target.Role.Name == entity.RoleSuperadmin && role.Name != entity.RoleSuperadmin {
		count, err := u.users.CountSuperadmins(ctx)
		if err != nil {
			return err
		}
		if count <= 1 {
			return entity.ErrLastSuperadmin
		}
	}
	if err := u.ensureEmailAvailable(ctx, email, target.ID); err != nil {
		return err
	}

	emailChanged := email != target.Email
	passwordHash := ""
	if input.Password != "" {
		if !entity.ValidPassword(input.Password) {
			return entity.ValidationError(map[string]string{"password": "validation.passwordLength"})
		}
		hash, err := u.hasher.Hash(input.Password)
		if err != nil {
			return err
		}
		passwordHash = hash
	}

	target.Name = name
	target.Email = email
	if emailChanged {
		target.EmailVerified = false
	}
	target.Status = status
	target.RoleID = role.ID
	target.PhoneNumber = clean(input.Phone)
	target.JobTitle = clean(input.JobTitle)
	target.Department = clean(input.Department)
	target.UpdatedAt = u.clock.Now().UTC()

	if err := u.users.UpdateWithPassword(ctx, target, passwordHash); err != nil {
		if errors.Is(err, port.ErrUniqueViolation) {
			return entity.ErrEmailInUse
		}
		return err
	}
	if status != entity.StatusActive {
		return u.sessions.DeleteByUserID(ctx, target.ID)
	}
	return nil
}

func (u *useCase) Delete(ctx context.Context, actor entity.Actor, id string) error {
	if !actor.Can(entity.PermissionManageUsers) {
		return entity.ErrForbidden
	}
	target, err := u.users.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrUserNotFound
		}
		return err
	}
	count, err := u.users.CountSuperadmins(ctx)
	if err != nil {
		return err
	}
	if !entity.CanDeleteUser(actor, target, count) {
		switch {
		case actor.ID == target.ID:
			return entity.ErrCannotDeleteSelf
		case target.Role != nil && target.Role.Name == entity.RoleSuperadmin:
			return entity.ErrLastSuperadmin
		default:
			return entity.ErrForbidden
		}
	}
	return u.users.Delete(ctx, target.ID)
}

func (u *useCase) Verify(ctx context.Context, actor entity.Actor, id string) error {
	if !actor.Can(entity.PermissionManageUsers) {
		return entity.ErrForbidden
	}
	target, err := u.users.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrUserNotFound
		}
		return err
	}
	if !entity.CanManageUser(actor, target) {
		return entity.ErrForbidden
	}
	if target.Status == entity.StatusPending {
		target.Status = entity.StatusActive
		target.UpdatedAt = u.clock.Now().UTC()
		return u.users.UpdateWithPassword(ctx, target, "")
	}
	return nil
}

func (u *useCase) SetRegistrationApproval(ctx context.Context, actor entity.Actor, enabled bool) error {
	if !actor.Can(entity.PermissionManageUsers) {
		return entity.ErrForbidden
	}
	value := "false"
	if enabled {
		value = "true"
	}
	return u.settings.Set(ctx, entity.SettingRegistrationApproval, value)
}

func (u *useCase) registrationApproval(ctx context.Context) (bool, error) {
	value, err := u.settings.Get(ctx, entity.SettingRegistrationApproval)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return entity.RegistrationApprovalEnabled(value), nil
}

func (u *useCase) assignableRole(ctx context.Context, actor entity.Actor, name string) (*entity.Role, error) {
	role, err := u.roles.GetByName(ctx, strings.TrimSpace(name))
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, entity.ErrRoleNotFound
		}
		return nil, err
	}
	if !entity.CanAssignRole(actor, role.Name) {
		return nil, entity.ErrForbidden
	}
	return role, nil
}

func (u *useCase) ensureEmailAvailable(ctx context.Context, email, excludeID string) error {
	existing, err := u.users.GetByEmail(ctx, email)
	if err == nil {
		if existing.ID != excludeID {
			return entity.ErrEmailInUse
		}
		return nil
	}
	if errors.Is(err, entity.ErrNotFound) {
		return nil
	}
	return err
}

func validate(input Input, requirePassword bool) (map[string]string, string, string) {
	fields := map[string]string{}
	name := strings.TrimSpace(input.Name)
	email := entity.NormalizeEmail(input.Email)
	if !entity.ValidateName(name) {
		fields["name"] = "validation.nameLength"
	}
	if !entity.ValidEmail(email) {
		fields["email"] = "validation.email"
	}
	if strings.TrimSpace(input.Role) == "" || !entity.ValidRoleName(strings.TrimSpace(input.Role)) {
		fields["role"] = "validation.invalidRole"
	}
	if input.Status != "" && !entity.IsUserStatus(input.Status) {
		fields["status"] = "validation.invalidStatus"
	}
	if requirePassword && !entity.ValidPassword(input.Password) {
		fields["password"] = "validation.passwordLength"
	}
	if !requirePassword && input.Password != "" && !entity.ValidPassword(input.Password) {
		fields["password"] = "validation.passwordLength"
	}
	return fields, name, email
}

func clean(value string) string { return strings.TrimSpace(value) }

func roleName(role *entity.Role) string {
	if role == nil {
		return ""
	}
	return role.Name
}
