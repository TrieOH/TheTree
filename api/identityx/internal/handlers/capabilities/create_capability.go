package capabilities

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) CreateCapability(ctx context.Context, req openapi.CreateCapabilityRequestObject) (openapi.CreateCapabilityResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	capability, err := h.ops.Create(ctx, models.CreateCapabilityInput{
		Resource:  req.Body.Resource,
		Action:    req.Body.Action,
		ProjectID: &req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateCapability201JSONResponse{
		Code: 201, Data: capability, Timestamp: time.Now(), Module: module,
	}, nil
}
