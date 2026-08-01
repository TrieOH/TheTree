package webhook_events

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) ListWebhookEvents(ctx context.Context, req openapi.ListWebhookEventsRequestObject) (openapi.ListWebhookEventsResponseObject, error) {
	events, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookEvents200JSONResponse{
		Code: 200, Data: &events, Timestamp: time.Now(), Module: module,
	}, nil
}
