package ports

import (
	"IdentityX/internal/shared/contracts"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type ProjectUserRepository interface {
	Register(ctx context.Context, toRegister contracts.ProjectUser) (*contracts.ProjectUser, error)
	GetByIDExternal(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) (*contracts.ProjectUser, error)
	GetByIDInternal(ctx context.Context, projectUserID, projectID uuid.UUID) (*contracts.ProjectUser, error)
	GetByEmailExternal(ctx context.Context, projectID uuid.UUID, email string, ownerID uuid.UUID) (*contracts.ProjectUser, error)
	GetByEmailInternal(ctx context.Context, projectID uuid.UUID, email string) (*contracts.ProjectUser, error)
	ListExternal(ctx context.Context, projectID, ownerID uuid.UUID) ([]contracts.ProjectUser, error)
	ListInternal(ctx context.Context, projectID uuid.UUID) ([]contracts.ProjectUser, error)
	Update(ctx context.Context, toUpdate contracts.ProjectUser, ownerID uuid.UUID) (*contracts.ProjectUser, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) error
	UpdateMetadata(ctx context.Context, userID, projectID uuid.UUID, metadata *json.RawMessage) error
	UpdateSubContext(ctx context.Context, userID, projectID uuid.UUID, subContext json.RawMessage) error
	Verify(ctx context.Context, userID uuid.UUID) (bool, error)
	BelongsToProject(ctx context.Context, userID, projectID uuid.UUID) (bool, error)
	ResetPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error
}
