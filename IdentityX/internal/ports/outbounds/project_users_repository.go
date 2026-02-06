package outbounds

import (
	"GoAuth/internal/domain/project_users"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type ProjectUserRepository interface {
	Register(ctx context.Context, toRegister project_users.ProjectUser) (*project_users.ProjectUser, error)
	GetByIDExternal(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) (*project_users.ProjectUser, error)
	GetByIDInternal(ctx context.Context, projectUserID, projectID uuid.UUID) (*project_users.ProjectUser, error)
	GetByEmailExternal(ctx context.Context, projectID uuid.UUID, email string, ownerID uuid.UUID) (*project_users.ProjectUser, error)
	GetByEmailInternal(ctx context.Context, projectID uuid.UUID, email string) (*project_users.ProjectUser, error)
	ListExternal(ctx context.Context, projectID, ownerID uuid.UUID) ([]project_users.ProjectUser, error)
	ListInternal(ctx context.Context, projectID uuid.UUID) ([]project_users.ProjectUser, error)
	Update(ctx context.Context, toUpdate project_users.ProjectUser, ownerID uuid.UUID) (*project_users.ProjectUser, error)
	Delete(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) error
	UpdateMetadata(ctx context.Context, userID, projectID uuid.UUID, metadata *json.RawMessage) error
	Verify(ctx context.Context, userID uuid.UUID) (bool, error)
	BelongsToProject(ctx context.Context, userID, projectID uuid.UUID) (bool, error)
}
