// Package capabilities implements the StrictServerInterface methods for
// the capabilities feature.
package capabilities

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
	"IdentityX/models"
)

const module = "IdentityX"

type Handlers struct {
	ops *services.Capabilities
}

func New(ops *services.Capabilities) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListCapabilities(ctx context.Context, req openapi.ListCapabilitiesRequestObject) (openapi.ListCapabilitiesResponseObject, error) {
	caps, err := h.ops.List(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListCapabilities200JSONResponse{
		Code: 200, Data: &caps, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateCapability(ctx context.Context, req openapi.CreateCapabilityRequestObject) (openapi.CreateCapabilityResponseObject, error) {
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
