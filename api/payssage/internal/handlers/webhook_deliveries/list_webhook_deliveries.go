package webhook_deliveries

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) ListWebhookDeliveries(ctx context.Context, req openapi.ListWebhookDeliveriesRequestObject) (openapi.ListWebhookDeliveriesResponseObject, error) {
	deliveries, err := h.ops.ListByEndpoint(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookDeliveries200JSONResponse{
		Code: 200, Data: &deliveries, Timestamp: time.Now(), Module: module,
	}, nil
}
