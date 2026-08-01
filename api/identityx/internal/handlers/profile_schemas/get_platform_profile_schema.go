package profile_schemas

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetPlatformProfileSchema(ctx context.Context, _ openapi.GetPlatformProfileSchemaRequestObject) (openapi.GetPlatformProfileSchemaResponseObject, error) {
	schema, err := h.ops.GetSchema(ctx, nil)
	if err != nil {
		return nil, err
	}
	return openapi.GetPlatformProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}
