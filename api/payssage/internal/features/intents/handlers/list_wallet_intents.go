package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) ListWalletIntents(ctx context.Context, req openapi.ListWalletIntentsRequestObject) (openapi.ListWalletIntentsResponseObject, error) {
	intents, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWalletIntents200JSONResponse{Items: intents}, nil
}
