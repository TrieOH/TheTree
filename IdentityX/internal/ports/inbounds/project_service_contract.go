package inbounds

import (
	"GoAuth/internal/domain/project"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectServiceInput struct {
	ProjectID   uuid.UUID
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

func OutputProjectSliceFromProjectSlice(src []project.Project) []OutputProject {
	dst := make([]OutputProject, 0, len(src))
	for _, p := range src {
		dst = append(dst, *OutputProjectFromProject(&p))
	}
	return dst
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

type ErrGeneratingProjectKeys struct {
	Cause error
}

func (e ErrGeneratingProjectKeys) Error() string {
	return "error generating project keys"
}

type ErrParsingProjectPublicKey struct {
	Cause error
}

func (e ErrParsingProjectPublicKey) Error() string {
	return "error parsing project public key"
}
