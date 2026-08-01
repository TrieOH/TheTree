// Package profile_schemas implements the StrictServerInterface methods for
// the profile_schemas feature.
package profile_schemas

import (
	"context"
	"encoding/json"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
	"IdentityX/models"
)

const module = "IdentityX"

type Handlers struct {
	ops *services.ProfileSchemas
}

func New(ops *services.ProfileSchemas) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) GetPlatformProfileSchema(ctx context.Context, _ openapi.GetPlatformProfileSchemaRequestObject) (openapi.GetPlatformProfileSchemaResponseObject, error) {
	schema, err := h.ops.GetSchema(ctx, nil)
	if err != nil {
		return nil, err
	}
	return openapi.GetPlatformProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) UpsertPlatformProfileSchema(ctx context.Context, req openapi.UpsertPlatformProfileSchemaRequestObject) (openapi.UpsertPlatformProfileSchemaResponseObject, error) {
	schema, err := h.ops.UpsertSchema(ctx, models.UpsertProfileSchemaInput{
		Schema: json.RawMessage(req.Body.Schema),
		Active: req.Body.Active,
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpsertPlatformProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetProjectProfileSchema(ctx context.Context, req openapi.GetProjectProfileSchemaRequestObject) (openapi.GetProjectProfileSchemaResponseObject, error) {
	schema, err := h.ops.GetSchema(ctx, &req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProjectProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) UpsertProjectProfileSchema(ctx context.Context, req openapi.UpsertProjectProfileSchemaRequestObject) (openapi.UpsertProjectProfileSchemaResponseObject, error) {
	schema, err := h.ops.UpsertSchema(ctx, models.UpsertProfileSchemaInput{
		ProjectID: &req.ProjectId,
		Schema:    json.RawMessage(req.Body.Schema),
		Active:    req.Body.Active,
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpsertProjectProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}
