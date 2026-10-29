package handlers

import (
	"context"

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
	return openapi.UpsertPlatformProfileSchema200JSONResponse(*schema), nil
}
