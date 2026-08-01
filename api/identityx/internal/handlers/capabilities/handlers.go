// Package capabilities implements the StrictServerInterface methods for
// the capabilities feature.
package capabilities

import (
	"context"
	"time"

	"IdentityX/internal/handler"
	"IdentityX/internal/services"
	"IdentityX/models"
)

const module = "IdentityX"

type Handlers struct {
	ops *services.Capabilities
}

func New(ops *services.Capabilities) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListCapabilities(ctx context.Context, req handler.ListCapabilitiesRequestObject) (handler.ListCapabilitiesResponseObject, error) {
	caps, err := h.ops.List(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return handler.ListCapabilities200JSONResponse{
		Code: 200, Data: &caps, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateCapability(ctx context.Context, req handler.CreateCapabilityRequestObject) (handler.CreateCapabilityResponseObject, error) {
	capability, err := h.ops.Create(ctx, models.CreateCapabilityInput{
		Resource:  req.Body.Resource,
		Action:    req.Body.Action,
		ProjectID: &req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return handler.CreateCapability201JSONResponse{
		Code: 201, Data: capability, Timestamp: time.Now(), Module: module,
	}, nil
}
