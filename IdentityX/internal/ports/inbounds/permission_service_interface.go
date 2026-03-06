package inbounds

import (
	"context"
)

type PermissionService interface {
	Create(ctx context.Context, in CreatePermissionInput) (*PermissionOutput, error)
	UpdateMeta(ctx context.Context, in UpdatePermissionInput) error
	Delete(ctx context.Context, in DeletePermissionInput) error

	// Read Operations //

	GetByIDExternal(ctx context.Context, in GetPermissionInput) (*PermissionOutput, error)
	ListByProject(ctx context.Context, in GetPermissionInput) ([]PermissionOutput, error)

	GiveDirect(ctx context.Context, in ManagePermissionInput) error
	TakeDirect(ctx context.Context, in ManagePermissionInput) error

	GetEffective(ctx context.Context, in ManagePermissionInput) ([]PermissionOutput, error)

	Check(ctx context.Context, in CheckPermissionInput) (bool, error)

	EnsureExists(ctx context.Context, in EnsurePermissionsInput) ([]EnsurePermissionResult, error)
}
