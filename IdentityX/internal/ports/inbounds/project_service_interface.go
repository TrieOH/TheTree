package inbounds

import (
	"context"

	"github.com/google/uuid"
)

type ProjectService interface {
	Create(ctx context.Context, in ProjectServiceInput) (*OutputProject, error)
	GetByID(ctx context.Context, projectID uuid.UUID) (*OutputProject, error)
	List(ctx context.Context) ([]OutputProject, error)
	GetJWKS(ctx context.Context, projectID uuid.UUID) (map[string]any, error)
	Update(ctx context.Context, in ProjectServiceInput) (*OutputProject, error)
	Delete(ctx context.Context, projectID uuid.UUID) error
	ListUsers(ctx context.Context, projectID uuid.UUID) ([]OutputProjectUser, error)
	GetUser(ctx context.Context, projectID, userID uuid.UUID) (*OutputProjectUser, error)
}
