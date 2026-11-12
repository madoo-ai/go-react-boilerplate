// Package roles implements the role and permission-matrix use cases.
package roles

import (
	"context"

	"go-react-boilerplate/internal/app/port"
	"go-react-boilerplate/internal/entity"
)

type Input struct {
	ID          string
	Name        string
	Label       string
	Description string
	Permissions []string
}

type ListResult struct {
	Roles       []entity.Role
	Permissions []entity.PermissionInfo
}

type UseCase interface {
	List(ctx context.Context, actor entity.Actor) (ListResult, error)
	Create(ctx context.Context, actor entity.Actor, input Input) error
	Update(ctx context.Context, actor entity.Actor, input Input) error
	Delete(ctx context.Context, actor entity.Actor, id string) error
}

type useCase struct {
	roles port.RoleStore
	clock port.Clock
}

func New(roles port.RoleStore, clock port.Clock) UseCase {
	return &useCase{roles: roles, clock: clock}
}
