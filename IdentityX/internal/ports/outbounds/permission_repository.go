package outbounds

import (
	"GoAuth/internal/domain/permissions"
	"context"

	"github.com/google/uuid"
)

type PermissionRepository interface {
	Create(ctx context.Context, toCreate permissions.Permission) (*permissions.Permission, error)

	// Read Operations //

	GetByIDInternal(ctx context.Context, id uuid.UUID) (*permissions.Permission, error)
	GetByIDExternal(ctx context.Context, id, projectID uuid.UUID) (*permissions.Permission, error)

	ListByProject(ctx context.Context, object, action *string, projectID uuid.UUID) ([]permissions.Permission, error)

	BelongsToProject(ctx context.Context, id, projectID uuid.UUID) (bool, error)
}
