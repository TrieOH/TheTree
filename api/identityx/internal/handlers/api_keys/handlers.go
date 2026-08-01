// Package api_keys implements the StrictServerInterface methods for the
// api_keys feature.
package api_keys

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
	"IdentityX/models"

	"github.com/google/uuid"
)

const module = "IdentityX"

func deref(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func derefSlice(v *[]openapi.UUID) []uuid.UUID {
	if v == nil {
		return nil
	}
	return *v
}

type Handlers struct {
	ops *services.APIKeys
}

func New(ops *services.APIKeys) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) CreateAPIKey(ctx context.Context, req openapi.CreateAPIKeyRequestObject) (openapi.CreateAPIKeyResponseObject, error) {
	key, raw, err := h.ops.Create(ctx, models.CreateAPIKeyInput{
		SubjectID:    req.Body.SubjectId,
		Capabilities: derefSlice(req.Body.Capabilities),
		Name:         deref(req.Body.Name),
		Env:          deref(req.Body.Env),
		ExpiresAt:    req.Body.ExpiresAt,
		ProjectID:    &req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	resp := models.CreateAPIKeyResponse{Key: key, RawKey: raw}
	return openapi.CreateAPIKey201JSONResponse{
		Code: 201, Data: &resp, Timestamp: time.Now(), Module: module,
	}, nil
}
