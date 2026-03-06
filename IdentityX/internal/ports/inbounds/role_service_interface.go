package inbounds

import (
	"context"
)

type RoleService interface {
	Create(ctx context.Context, in RoleInput) (*RoleOutput, error)
	UpdateDescription(ctx context.Context, in RoleInput) error
	UpdateMeta(ctx context.Context, in RoleInput) error
	Delete(ctx context.Context, in RoleInput) error
	GetByIDExternal(ctx context.Context, in GetRoleInput) (*RoleOutput, error)
	GetByName(ctx context.Context, in GetRoleInput) (*RoleOutput, error)
	ListByProject(ctx context.Context, in GetRoleInput) ([]RoleOutput, error)

	AddPermission(ctx context.Context, in RolePermissionInput) error
	RemovePermission(ctx context.Context, in RolePermissionInput) error

	GetPermissions(ctx context.Context, in RolePermissionInput) ([]PermissionOutput, error)

	GiveRole(ctx context.Context, in ManageRoleInput) error
	TakeRole(ctx context.Context, in ManageRoleInput) error

	GiveRoleByName(ctx context.Context, in ManageRoleInput) error
	TakeRoleByName(ctx context.Context, in ManageRoleInput) error

	GetUserRoles(ctx context.Context, in GetRoleInput) ([]RoleOutput, error)

	EnsureExists(ctx context.Context, in EnsureRolesInput) ([]EnsureRoleResult, error)
}
