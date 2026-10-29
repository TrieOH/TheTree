package handlers

import (
	"context"

	"IdentityX/internal/features/profile_schemas/service"
	"IdentityX/internal/openapi"
	"IdentityX/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) GetPlatformProfileSchema(ctx context.Context, _ openapi.GetPlatformProfileSchemaRequestObject) (openapi.GetPlatformProfileSchemaResponseObject, error) {
	schema, err := h.ops.GetSchema(ctx, nil)
	if err != nil {
		return nil, err
	}
	return openapi.GetPlatformProfileSchema200JSONResponse(*schema), nil
}

func (h *Handlers) GetProjectProfileSchema(ctx context.Context, req openapi.GetProjectProfileSchemaRequestObject) (openapi.GetProjectProfileSchemaResponseObject, error) {
	schema, err := h.ops.GetSchema(ctx, &req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProjectProfileSchema200JSONResponse(*schema), nil
}

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

func (h *Handlers) UpsertProjectProfileSchema(ctx context.Context, req openapi.UpsertProjectProfileSchemaRequestObject) (openapi.UpsertProjectProfileSchemaResponseObject, error) {
	schema, err := h.ops.UpsertSchema(ctx, models.UpsertProfileSchemaInput{
		ProjectID: &req.ProjectId,
		Schema:    req.Body.Schema,
		Active:    req.Body.Active,
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpsertProjectProfileSchema200JSONResponse(*schema), nil
}
