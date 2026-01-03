package inbounds

import (
	"GoAuth/internal/domain/project"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectService interface {
	Create(ctx context.Context, in CreateProjectInput) (*OutputProject, error)
	GetByID(ctx context.Context, projectID string) (*OutputProject, error)
	List(ctx context.Context) ([]OutputProject, error)
	GetJWKS(ctx context.Context, projectID string) (map[string]any, error)
	Update(ctx context.Context, in UpdateProjectInput) (*OutputProject, error)
	Delete(ctx context.Context, projectID string) error
}

type CreateProjectInput struct {
	ProjectName string
	Metadata    json.RawMessage
}

type UpdateProjectInput struct {
	ProjectID   string
	ProjectName string
	Metadata    json.RawMessage
}

type OutputProject struct {
	ID          uuid.UUID
	ProjectName string
	OwnerID     uuid.UUID
	Metadata    json.RawMessage
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func OutputProjectFromProject(p *project.Project) *OutputProject {
	return &OutputProject{
		ID:          p.ID,
		ProjectName: p.ProjectName,
		OwnerID:     p.OwnerID,
		Metadata:    p.Metadata,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func OutputProjectSliceFromProjectSlice(src []project.Project) []OutputProject {
	dst := make([]OutputProject, 0, len(src))
	for _, p := range src {
		dst = append(dst, *OutputProjectFromProject(&p))
	}
	return dst
}
