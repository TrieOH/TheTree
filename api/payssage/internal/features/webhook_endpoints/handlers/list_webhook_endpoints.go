package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) ListWebhookEndpoints(ctx context.Context, req openapi.ListWebhookEndpointsRequestObject) (openapi.ListWebhookEndpointsResponseObject, error) {
	endpoints, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWebhookEndpoints200JSONResponse{Items: endpoints}, nil
}
