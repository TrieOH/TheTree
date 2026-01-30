package inbounds

import (
	"context"
)

type PermissionService interface {
	Create(ctx context.Context, in CreatePermissionInput) (*PermissionOutput, error)

	// Read Operations //

	GetByIDExternal(ctx context.Context, in GetPermissionInput) (*PermissionOutput, error)
	ListByProject(ctx context.Context, in GetPermissionInput) ([]PermissionOutput, error)

	GiveDirect(ctx context.Context, in ManagePermissionInput) error
	TakeDirect(ctx context.Context, in ManagePermissionInput) error

	GetEffective(ctx context.Context, in ManagePermissionInput) ([]PermissionOutput, error)

	Check(ctx context.Context, in CheckPermissionInput) (bool, error)
}
