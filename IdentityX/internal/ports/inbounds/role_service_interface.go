package inbounds

import (
	"context"
)

type RoleService interface {
	Create(ctx context.Context, in RoleInput) (*RoleOutput, error)
	UpdateDescription(ctx context.Context, in RoleInput) error
	GetByIDExternal(ctx context.Context, in GetRoleInput) (*RoleOutput, error)
	GetByName(ctx context.Context, in GetRoleInput) (*RoleOutput, error)
	ListByProject(ctx context.Context, in GetRoleInput) ([]RoleOutput, error)
}
