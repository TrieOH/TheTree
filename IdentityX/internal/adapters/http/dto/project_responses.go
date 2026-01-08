package dto

import (
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
	"time"
)

type ProjectResponse struct {
	ID          string          `json:"id"`
	ProjectName string          `json:"project_name"`
	OwnerID     string          `json:"owner_id"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func ProjectSliceToProjectResponseSlice(src []inbounds.OutputProject) []ProjectResponse {
	dst := make([]ProjectResponse, 0, len(src))
	for _, p := range src {
		dst = append(dst, ProjectToResponse(&p))
	}
	return dst
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
