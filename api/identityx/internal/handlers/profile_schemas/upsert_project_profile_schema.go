package profile_schemas

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) UpsertProjectProfileSchema(ctx context.Context, req openapi.UpsertProjectProfileSchemaRequestObject) (openapi.UpsertProjectProfileSchemaResponseObject, error) {
	schema, err := h.ops.UpsertSchema(ctx, models.UpsertProfileSchemaInput{
		ProjectID: &req.ProjectId,
		Schema:    req.Body.Schema,
		Active:    req.Body.Active,
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpsertProjectProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}
