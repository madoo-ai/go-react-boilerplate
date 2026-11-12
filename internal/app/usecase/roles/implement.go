package roles

import (
	"context"
	"errors"
	"strings"

	"go-react-boilerplate/internal/app/port"
	"go-react-boilerplate/internal/entity"
)

func (u *useCase) List(ctx context.Context, actor entity.Actor) (ListResult, error) {
	if !actor.Can(entity.PermissionManageRoles) {
		return ListResult{}, entity.ErrForbidden
	}
	roles, err := u.roles.List(ctx)
	if err != nil {
		return ListResult{}, err
	}
	permissions, err := u.roles.ListPermissions(ctx)
	if err != nil {
		return ListResult{}, err
	}
	return ListResult{Roles: roles, Permissions: permissions}, nil
}

func (u *useCase) Create(ctx context.Context, actor entity.Actor, input Input) error {
	if !actor.Can(entity.PermissionManageRoles) {
		return entity.ErrForbidden
	}
	fields, name, label, description := validate(input)
	if len(fields) > 0 {
		return entity.ValidationError(fields)
	}
	if _, err := u.roles.GetByName(ctx, name); err == nil {
		return entity.ErrRoleInUse
	} else if !errors.Is(err, entity.ErrNotFound) {
		return err
	}
	now := u.clock.Now().UTC()
	role := &entity.Role{
		ID: entity.NewUUID(), Name: name, Label: label, Description: description,
		IsSystem: false, Permissions: validPermissions(input.Permissions),
		CreatedAt: now, UpdatedAt: now,
	}
	if err := u.roles.Create(ctx, role); err != nil {
		if errors.Is(err, port.ErrUniqueViolation) {
			return entity.ErrRoleInUse
		}
		return err
	}
	return nil
}

func (u *useCase) Update(ctx context.Context, actor entity.Actor, input Input) error {
	if !actor.Can(entity.PermissionManageRoles) {
		return entity.ErrForbidden
	}
	role, err := u.roles.GetByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrRoleNotFound
		}
		return err
	}
	if role.Name == entity.RoleSuperadmin {
		return entity.ErrImmutableRole
	}
	fields, name, label, description := validate(input)
	if len(fields) > 0 {
		return entity.ValidationError(fields)
	}
	if role.Name == entity.RoleUser && name != entity.RoleUser {
		return entity.ErrDefaultRole
	}
	duplicate, err := u.roles.GetByName(ctx, name)
	if err == nil && duplicate.ID != role.ID {
		return entity.ErrRoleInUse
	}
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return err
	}
	role.Name = name
	role.Label = label
	role.Description = description
	role.Permissions = validPermissions(input.Permissions)
	role.UpdatedAt = u.clock.Now().UTC()
	if err := u.roles.Update(ctx, role); err != nil {
		if errors.Is(err, port.ErrUniqueViolation) {
			return entity.ErrRoleInUse
		}
		return err
	}
	return nil
}

func (u *useCase) Delete(ctx context.Context, actor entity.Actor, id string) error {
	if !actor.Can(entity.PermissionManageRoles) {
		return entity.ErrForbidden
	}
	role, err := u.roles.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrRoleNotFound
		}
		return err
	}
	if role.Name == entity.RoleSuperadmin {
		return entity.ErrImmutableRole
	}
	if role.Name == entity.RoleUser {
		return entity.ErrDefaultRole
	}
	if role.UserCount > 0 {
		return entity.ErrRoleInUse
	}
	return u.roles.Delete(ctx, role.ID)
}

func validate(input Input) (map[string]string, string, string, string) {
	fields := map[string]string{}
	name := strings.TrimSpace(input.Name)
	label := strings.TrimSpace(input.Label)
	description := strings.TrimSpace(input.Description)
	if !entity.ValidRoleName(name) {
		fields["name"] = "validation.roleName"
	}
	if len([]rune(label)) < 1 || len([]rune(label)) > 80 {
		fields["label"] = "validation.required"
	}
	if len([]rune(description)) > 240 {
		fields["description"] = "validation.max"
	}
	return fields, name, label, description
}

func validPermissions(values []string) []entity.Permission {
	out := make([]entity.Permission, 0, len(values))
	seen := map[entity.Permission]bool{}
	for _, value := range values {
		if !entity.IsPermission(value) {
			continue
		}
		permission := entity.Permission(value)
		if !seen[permission] {
			seen[permission] = true
			out = append(out, permission)
		}
	}
	return out
}
