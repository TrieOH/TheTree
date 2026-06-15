package ports

import (
	"context"

	"payssage/models"

	"github.com/google/uuid"
)

type WorkspaceRepo interface {
	Create(ctx context.Context, toCreate models.Workspace) (*models.Workspace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Workspace, error)
	GetByName(ctx context.Context, name string, userID uuid.UUID) (*models.Workspace, error)
	List(ctx context.Context, userID uuid.UUID) ([]models.Workspace, error)
	EnableSandbox(ctx context.Context, id uuid.UUID) (*models.Workspace, error)
	DisableSandbox(ctx context.Context, id uuid.UUID) (*models.Workspace, error)
}
