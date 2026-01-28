package outbounds

import (
	"GoAuth/internal/domain/permissions"
	"GoAuth/internal/domain/roles"
	"context"

	"github.com/google/uuid"
)

type RoleRepository interface {
	Create(ctx context.Context, toCreate roles.Role) (*roles.Role, error)
	UpdateDescription(ctx context.Context, description string, id uuid.UUID, projectID *uuid.UUID) error
	GetByIDInternal(ctx context.Context, id uuid.UUID) (*roles.Role, error)
	GetByIDExternal(ctx context.Context, id, projectID uuid.UUID) (*roles.Role, error)
	GetByName(ctx context.Context, name string, projectID *uuid.UUID) (*roles.Role, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]roles.Role, error)

	BelongsToProject(ctx context.Context, id, projectID uuid.UUID) (bool, error)

	AddPermission(ctx context.Context, id uuid.UUID, permissionID uuid.UUID) error
	RemovePermission(ctx context.Context, id uuid.UUID, permissionID uuid.UUID) error

	GetPermissions(ctx context.Context, id, projectID uuid.UUID) ([]permissions.Permission, error)

	GiveRole(ctx context.Context, id, identityID uuid.UUID, scopeID *uuid.UUID) error
	TakeRole(ctx context.Context, id, identityID uuid.UUID, scopeID *uuid.UUID) error

	GetUserRoles(ctx context.Context, identityID, projectID uuid.UUID) ([]roles.Role, error)
}
