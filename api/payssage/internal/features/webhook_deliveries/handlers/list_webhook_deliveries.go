package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) ListWebhookDeliveries(ctx context.Context, req openapi.ListWebhookDeliveriesRequestObject) (openapi.ListWebhookDeliveriesResponseObject, error) {
	deliveries, err := h.ops.ListByEndpoint(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookDeliveries200JSONResponse{Items: deliveries}, nil
}
