package dto

import (
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RoleResponse struct {
	ID          uuid.UUID        `json:"id"`
	ProjectID   *uuid.UUID       `json:"project_id"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Meta        *json.RawMessage `json:"meta"`

	ScopeID    *uuid.UUID `json:"scope_id"`
	ScopeName  *string    `json:"scope_name"`
	ExternalID *string    `json:"external_id"`
}

func RoleOutputToRoleResponse(role inbounds.RoleOutput) *RoleResponse {
	return &RoleResponse{
		ID:          role.Role.ID,
		ProjectID:   role.Role.ProjectID,
		Name:        role.Role.Name,
		Description: role.Role.Description,
		CreatedAt:   role.Role.CreatedAt,
		UpdatedAt:   role.Role.UpdatedAt,
		ScopeID:     role.Role.ScopeID,
		ScopeName:   role.Role.ScopeName,
		Meta:        role.Role.Meta,
		ExternalID:  role.Role.ExternalID,
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
			ScopeID:     role.Role.ScopeID,
			Meta:        role.Role.Meta,
			ScopeName:   role.Role.ScopeName,
			ExternalID:  role.Role.ExternalID,
		})
	}
	return out
}
