package inbounds

import (
	"context"
)

type PermissionService interface {
	Create(ctx context.Context, in CreatePermissionInput) (*PermissionOutput, error)

	// Read Operations //

	GetByIDExternal(ctx context.Context, in GetPermissionInput) (*PermissionOutput, error)
	ListByProject(ctx context.Context, in GetPermissionInput) ([]PermissionOutput, error)
}
