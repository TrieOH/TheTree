package handlers

import (
	"context"

	"IdentityX/internal/features/api_keys/service"
	"IdentityX/internal/openapi"
	"IdentityX/models"

	"github.com/google/uuid"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

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
	return openapi.CreateAPIKey201JSONResponse(resp), nil
}
