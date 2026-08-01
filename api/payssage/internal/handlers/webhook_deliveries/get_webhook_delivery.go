package webhook_deliveries

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) GetWebhookDelivery(ctx context.Context, req openapi.GetWebhookDeliveryRequestObject) (openapi.GetWebhookDeliveryResponseObject, error) {
	delivery, err := h.ops.GetByID(ctx, req.DeliveryId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWebhookDelivery200JSONResponse{
		Code: 200, Data: delivery, Timestamp: time.Now(), Module: module,
	}, nil
}
