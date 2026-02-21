package dto

import (
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ScopeResponse struct {
	ID         uuid.UUID        `json:"id"`
	ParentID   *uuid.UUID       `json:"parent_id"`
	ProjectID  *uuid.UUID       `json:"project_id"`
	ExternalID *string          `json:"external_id"`
	Type       string           `json:"type"`
	Name       *string          `json:"name"`
	Meta       *json.RawMessage `json:"meta"`
	CreatedAt  time.Time        `json:"created_at"`
}

func ScopeOutputToScopeResponse(in *inbounds.ScopeOutput) ScopeResponse {
	return ScopeResponse{
		ID:         in.ID,
		ParentID:   in.ParentID,
		ProjectID:  in.ProjectID,
		ExternalID: in.ExternalID,
		Type:       string(in.Type),
		Name:       in.Name,
		Meta:       in.Meta,
		CreatedAt:  in.CreatedAt,
	}
}

func ScopeOutputSliceToScopeResponseSlice(in []inbounds.ScopeOutput) []ScopeResponse {
	if in == nil {
		return nil
	}

	out := make([]ScopeResponse, 0, len(in))
	for _, i := range in {
		out = append(out, ScopeOutputToScopeResponse(&i))
	}
	return out
}
