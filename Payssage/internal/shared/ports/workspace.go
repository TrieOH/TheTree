package ports

import (
	"context"
	"payssage/internal/shared/contracts"

	"github.com/google/uuid"
)

type WorkspaceRepo interface {
	Create(ctx context.Context, toCreate contracts.Workspace) (*contracts.Workspace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Workspace, error)
	GetByName(ctx context.Context, name string, userID uuid.UUID) (*contracts.Workspace, error)
	List(ctx context.Context, userID uuid.UUID) ([]contracts.Workspace, error)
	EnableSandbox(ctx context.Context, id uuid.UUID) (*contracts.Workspace, error)
	DisableSandbox(ctx context.Context, id uuid.UUID) (*contracts.Workspace, error)
}
