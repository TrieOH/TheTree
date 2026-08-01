package profile_schemas

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetProjectProfileSchema(ctx context.Context, req openapi.GetProjectProfileSchemaRequestObject) (openapi.GetProjectProfileSchemaResponseObject, error) {
	schema, err := h.ops.GetSchema(ctx, &req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProjectProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}
