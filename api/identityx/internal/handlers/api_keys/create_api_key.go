package api_keys

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) CreateAPIKey(ctx context.Context, req openapi.CreateAPIKeyRequestObject) (openapi.CreateAPIKeyResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

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
