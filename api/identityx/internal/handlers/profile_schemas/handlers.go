// Package profile_schemas implements the StrictServerInterface methods for
// the profile_schemas feature.
package profile_schemas

import (
	"context"
	"encoding/json"
	"time"

	"IdentityX/internal/handler"
	"IdentityX/internal/services"
	"IdentityX/models"
)

const module = "IdentityX"

type Handlers struct {
	ops *services.ProfileSchemas
}

func New(ops *services.ProfileSchemas) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) GetPlatformProfileSchema(ctx context.Context, _ handler.GetPlatformProfileSchemaRequestObject) (handler.GetPlatformProfileSchemaResponseObject, error) {
	schema, err := h.ops.GetSchema(ctx, nil)
	if err != nil {
		return nil, err
	}
	return handler.GetPlatformProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) UpsertPlatformProfileSchema(ctx context.Context, req handler.UpsertPlatformProfileSchemaRequestObject) (handler.UpsertPlatformProfileSchemaResponseObject, error) {
	schema, err := h.ops.UpsertSchema(ctx, models.UpsertProfileSchemaInput{
		Schema: json.RawMessage(req.Body.Schema),
		Active: req.Body.Active,
	})
	if err != nil {
		return nil, err
	}
	return handler.UpsertPlatformProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetProjectProfileSchema(ctx context.Context, req handler.GetProjectProfileSchemaRequestObject) (handler.GetProjectProfileSchemaResponseObject, error) {
	schema, err := h.ops.GetSchema(ctx, &req.ProjectId)
	if err != nil {
		return nil, err
	}
	return handler.GetProjectProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) UpsertProjectProfileSchema(ctx context.Context, req handler.UpsertProjectProfileSchemaRequestObject) (handler.UpsertProjectProfileSchemaResponseObject, error) {
	schema, err := h.ops.UpsertSchema(ctx, models.UpsertProfileSchemaInput{
		ProjectID: &req.ProjectId,
		Schema:    json.RawMessage(req.Body.Schema),
		Active:    req.Body.Active,
	})
	if err != nil {
		return nil, err
	}
	return handler.UpsertProjectProfileSchema200JSONResponse{
		Code: 200, Data: schema, Timestamp: time.Now(), Module: module,
	}, nil
}
