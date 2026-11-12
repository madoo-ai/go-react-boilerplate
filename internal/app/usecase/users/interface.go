// Package users implements the account-management use cases. Every mutation
// re-checks the actor's permission and the target's protection rules close to
// the data.
package users

import (
	"context"

	"go-react-boilerplate/internal/app/port"
	"go-react-boilerplate/internal/entity"
)

type Input struct {
	ID         string
	Name       string
	Email      string
	Phone      string
	JobTitle   string
	Department string
	Role       string
	Status     string
	Password   string
}

type UserListItem struct {
	User      entity.User
	CanManage bool
	CanDelete bool
}

type ListResult struct {
	Users                []UserListItem
	Total                int
	Page                 int
	PageCount            int
	PageSize             int
	Roles                []entity.Role
	AssignableRoles      []entity.Role
	RegistrationApproval bool
	SuperadminCount      int
}

type EditResult struct {
	User            entity.User
	AssignableRoles []entity.Role
	CanManage       bool
}

type UseCase interface {
	List(ctx context.Context, actor entity.Actor, filters entity.UserFilters) (ListResult, error)
	GetForEdit(ctx context.Context, actor entity.Actor, id string) (EditResult, error)
	Create(ctx context.Context, actor entity.Actor, input Input) error
	Update(ctx context.Context, actor entity.Actor, input Input) error
	Delete(ctx context.Context, actor entity.Actor, id string) error
	Verify(ctx context.Context, actor entity.Actor, id string) error
	SetRegistrationApproval(ctx context.Context, actor entity.Actor, enabled bool) error
}

type useCase struct {
	users    port.UserStore
	roles    port.RoleStore
	settings port.SettingStore
	sessions port.SessionStore
	hasher   port.PasswordHasher
	clock    port.Clock
}

func New(users port.UserStore, roles port.RoleStore, settings port.SettingStore, sessions port.SessionStore, hasher port.PasswordHasher, clock port.Clock) UseCase {
	return &useCase{users: users, roles: roles, settings: settings, sessions: sessions, hasher: hasher, clock: clock}
}
