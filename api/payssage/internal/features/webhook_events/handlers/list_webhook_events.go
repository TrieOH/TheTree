package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) ListWebhookEvents(ctx context.Context, req openapi.ListWebhookEventsRequestObject) (openapi.ListWebhookEventsResponseObject, error) {
	events, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookEvents200JSONResponse{Items: events}, nil
}
