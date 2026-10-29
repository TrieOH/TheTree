package handlers

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) CreateCapability(ctx context.Context, req openapi.CreateCapabilityRequestObject) (openapi.CreateCapabilityResponseObject, error) {
	capability, err := h.ops.Create(ctx, models.CreateCapabilityInput{
		Resource:  req.Body.Resource,
		Action:    req.Body.Action,
		ProjectID: &req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateCapability201JSONResponse(*capability), nil
}
