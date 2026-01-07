package inbounds

import (
	"context"
)

type ProjectService interface {
	Create(ctx context.Context, in ProjectServiceInput) (*OutputProject, error)
	GetByID(ctx context.Context, projectID string) (*OutputProject, error)
	List(ctx context.Context) ([]OutputProject, error)
	GetJWKS(ctx context.Context, projectID string) (map[string]any, error)
	Update(ctx context.Context, in ProjectServiceInput) (*OutputProject, error)
	Delete(ctx context.Context, projectID string) error
}
