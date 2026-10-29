package collectors

import (
	"context"

	"payssage/internal/features/collectors/service"
	"payssage/internal/openapi"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) GetCollector(ctx context.Context, req openapi.GetCollectorRequestObject) (openapi.GetCollectorResponseObject, error) {
	collector, err := h.ops.GetByID(ctx, req.CollectorId)
	if err != nil {
		return nil, err
	}
	return openapi.GetCollector200JSONResponse(*collector), nil
}

func (h *Handlers) ListCollectors(ctx context.Context, _ openapi.ListCollectorsRequestObject) (openapi.ListCollectorsResponseObject, error) {
	collectors, err := h.ops.ListOwned(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListCollectors200JSONResponse{Items: collectors}, nil
}

func (h *Handlers) ListOrganizationCollectors(ctx context.Context, req openapi.ListOrganizationCollectorsRequestObject) (openapi.ListOrganizationCollectorsResponseObject, error) {
	collectors, err := h.ops.ListByOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationCollectors200JSONResponse{Items: collectors}, nil
}
