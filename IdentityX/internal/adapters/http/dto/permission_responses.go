package dto

import (
	"GoAuth/internal/ports/inbounds"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PermissionResponse struct {
	ID         uuid.UUID        `json:"id"`
	ProjectID  *uuid.UUID       `json:"project_id"`
	Object     string           `json:"object"`
	Action     string           `json:"action"`
	Conditions *json.RawMessage `json:"conditions"`
	CreatedAt  time.Time        `json:"created_at"`
}

func PermissionOutputToPermissionResponse(in inbounds.PermissionOutput) PermissionResponse {
	return PermissionResponse{
		ID:         in.Permission.ID,
		ProjectID:  in.Permission.ProjectID,
		Object:     in.Permission.Object,
		Action:     in.Permission.Action,
		Conditions: in.Permission.Conditions,
		CreatedAt:  in.Permission.CreatedAt,
	}
}

func PermissionOutputSliceToPermissionResponseSlice(in []inbounds.PermissionOutput) []PermissionResponse {
	if in == nil {
		return nil
	}

	out := make([]PermissionResponse, 0, len(in))
	for _, i := range in {
		out = append(out, PermissionOutputToPermissionResponse(i))
	}
	return out
}
