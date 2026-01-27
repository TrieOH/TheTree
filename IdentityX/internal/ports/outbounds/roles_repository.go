package outbounds

import (
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
}
