package capabilities

import (
	"context"

	"IdentityX/internal/features/capabilities/service"
	"IdentityX/internal/openapi"
	"IdentityX/models"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

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

func (h *Handlers) ListCapabilities(ctx context.Context, req openapi.ListCapabilitiesRequestObject) (openapi.ListCapabilitiesResponseObject, error) {
	caps, err := h.ops.List(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListCapabilities200JSONResponse{Items: caps}, nil
}
