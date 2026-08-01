// Package webhook_deliveries implements the StrictServerInterface methods
// for the webhook_deliveries feature.
package webhook_deliveries

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/internal/services"
)

const module = "Payssage"

type Handlers struct {
	ops *services.WebhookDeliveries
}

func New(ops *services.WebhookDeliveries) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListWebhookDeliveries(ctx context.Context, req openapi.ListWebhookDeliveriesRequestObject) (openapi.ListWebhookDeliveriesResponseObject, error) {
	deliveries, err := h.ops.ListByEndpoint(ctx, req.EndpointId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookDeliveries200JSONResponse{
		Code: 200, Data: &deliveries, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetWebhookDelivery(ctx context.Context, req openapi.GetWebhookDeliveryRequestObject) (openapi.GetWebhookDeliveryResponseObject, error) {
	delivery, err := h.ops.GetByID(ctx, req.DeliveryId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookDelivery200JSONResponse{
		Code: 200, Data: delivery, Timestamp: time.Now(), Module: module,
	}, nil
}
