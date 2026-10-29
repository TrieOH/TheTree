package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetPlatformProfileSchema(ctx context.Context, _ openapi.GetPlatformProfileSchemaRequestObject) (openapi.GetPlatformProfileSchemaResponseObject, error) {
	schema, err := h.ops.GetSchema(ctx, nil)
	if err != nil {
		return nil, err
	}
	return openapi.GetPlatformProfileSchema200JSONResponse(*schema), nil
}
