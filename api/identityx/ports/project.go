package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type ProjectRepo interface {
	Create(ctx context.Context, project models.Project) (*models.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	// ListAll returns every project, platform or tenant — the Key-lifecycle
	// module's EnsureAll crosses it to provision every scope.
	ListAll(ctx context.Context) ([]models.Project, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]models.Project, error)
	ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Project, error)
	ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Project, error)
	AddMember(ctx context.Context, toCreate models.ProjectMember) error
	RemoveMember(ctx context.Context, actorID, projectID uuid.UUID) error
	GetMember(ctx context.Context, actorID, projectID uuid.UUID) (*models.ProjectMember, error)
	GetRole(ctx context.Context, actorID, projectID uuid.UUID) (models.ProjectRole, error)
	ListMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error)
}
