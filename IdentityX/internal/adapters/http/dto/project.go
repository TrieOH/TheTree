package dto

import (
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
	"time"
)

type CreateProjectRequest struct {
	ProjectName string          `json:"project_name" validate:"required,max=255"`
	Metadata    json.RawMessage `json:"metadata"`
}

type UpdateProjectRequest struct {
	ProjectName string          `json:"project_name" validate:"required,max=255"`
	Metadata    json.RawMessage `json:"metadata"`
}

type ProjectResponse struct {
	ID          string          `json:"id"`
	ProjectName string          `json:"project_name"`
	OwnerID     string          `json:"owner_id"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func ProjectToResponse(r *inbounds.OutputProject) ProjectResponse {
	return ProjectResponse{
		ID:          r.ID.String(),
		ProjectName: r.ProjectName,
		OwnerID:     r.OwnerID.String(),
		Metadata:    r.Metadata,
		IsActive:    r.IsActive,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func ProjectSliceToProjectResponseSlice(src []inbounds.OutputProject) []ProjectResponse {
	dst := make([]ProjectResponse, 0, len(src))
	for _, p := range src {
		dst = append(dst, ProjectToResponse(&p))
	}
	return dst
}
