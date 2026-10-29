package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) GetWebhookDelivery(ctx context.Context, req openapi.GetWebhookDeliveryRequestObject) (openapi.GetWebhookDeliveryResponseObject, error) {
	delivery, err := h.ops.GetByID(ctx, req.DeliveryId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookDelivery200JSONResponse(*delivery), nil
}
