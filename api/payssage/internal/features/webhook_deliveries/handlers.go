package webhook_deliveries

import (
	"context"

	"payssage/internal/features/webhook_deliveries/service"
	"payssage/internal/openapi"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) GetWebhookDelivery(ctx context.Context, req openapi.GetWebhookDeliveryRequestObject) (openapi.GetWebhookDeliveryResponseObject, error) {
	delivery, err := h.ops.GetByID(ctx, req.DeliveryId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookDelivery200JSONResponse(*delivery), nil
}

func (h *Handlers) ListWebhookDeliveries(ctx context.Context, req openapi.ListWebhookDeliveriesRequestObject) (openapi.ListWebhookDeliveriesResponseObject, error) {
	deliveries, err := h.ops.ListByEndpoint(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookDeliveries200JSONResponse{Items: deliveries}, nil
}
