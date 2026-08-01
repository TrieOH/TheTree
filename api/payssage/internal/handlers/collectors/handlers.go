// Package collectors implements the StrictServerInterface methods for the
// collectors feature.
package collectors

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/internal/services"
)

const module = "Payssage"

type Handlers struct {
	ops *services.Collectors
}

func New(ops *services.Collectors) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListCollectors(ctx context.Context, _ openapi.ListCollectorsRequestObject) (openapi.ListCollectorsResponseObject, error) {
	collectors, err := h.ops.ListOwned(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListCollectors200JSONResponse{
		Code: 200, Data: &collectors, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetCollector(ctx context.Context, req openapi.GetCollectorRequestObject) (openapi.GetCollectorResponseObject, error) {
	collector, err := h.ops.GetByID(ctx, req.CollectorId)
	if err != nil {
		return nil, err
	}
	return openapi.GetCollector200JSONResponse{
		Code: 200, Data: collector, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListOrganizationCollectors(ctx context.Context, req openapi.ListOrganizationCollectorsRequestObject) (openapi.ListOrganizationCollectorsResponseObject, error) {
	collectors, err := h.ops.ListByOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationCollectors200JSONResponse{
		Code: 200, Data: &collectors, Timestamp: time.Now(), Module: module,
	}, nil
}
