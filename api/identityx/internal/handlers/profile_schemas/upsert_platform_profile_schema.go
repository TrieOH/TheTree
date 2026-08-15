package profile_schemas

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) UpsertPlatformProfileSchema(ctx context.Context, req openapi.UpsertPlatformProfileSchemaRequestObject) (openapi.UpsertPlatformProfileSchemaResponseObject, error) {
	schema, err := h.ops.UpsertSchema(ctx, models.UpsertProfileSchemaInput{
		Schema: req.Body.Schema,
		Active: req.Body.Active,
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpsertPlatformProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}
