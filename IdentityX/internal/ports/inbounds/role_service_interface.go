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

	AddPermission(ctx context.Context, in RolePermissionInput) error
	RemovePermission(ctx context.Context, in RolePermissionInput) error

	GetPermissions(ctx context.Context, in RolePermissionInput) ([]PermissionOutput, error)

	GiveRole(ctx context.Context, in ManageRoleInput) error
	TakeRole(ctx context.Context, in ManageRoleInput) error

	GetUserRoles(ctx context.Context, in GetRoleInput) ([]RoleOutput, error)
}
