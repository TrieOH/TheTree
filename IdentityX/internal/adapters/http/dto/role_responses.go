package dto

import (
	"GoAuth/internal/ports/inbounds"
	"time"

	"github.com/google/uuid"
)

type RoleResponse struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   *uuid.UUID `json:"project_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func RoleOutputToRoleResponse(role inbounds.RoleOutput) *RoleResponse {
	return &RoleResponse{
		ID:          role.Role.ID,
		ProjectID:   role.Role.ProjectID,
		Name:        role.Role.Name,
		Description: role.Role.Description,
		CreatedAt:   role.Role.CreatedAt,
		UpdatedAt:   role.Role.UpdatedAt,
	}
}

func RoleOutputSliceToRoleResponseSlice(in []inbounds.RoleOutput) []RoleResponse {
	if in == nil {
		return nil
	}

	out := make([]RoleResponse, 0, len(in))
	for _, role := range in {
		out = append(out, RoleResponse{
			ID:          role.Role.ID,
			ProjectID:   role.Role.ProjectID,
			Name:        role.Role.Name,
			Description: role.Role.Description,
			CreatedAt:   role.Role.CreatedAt,
			UpdatedAt:   role.Role.UpdatedAt,
		})
	}
	return out
}
