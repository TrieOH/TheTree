package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) ListWalletSellers(ctx context.Context, req openapi.ListWalletSellersRequestObject) (openapi.ListWalletSellersResponseObject, error) {
	sellers, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWalletSellers200JSONResponse{Items: sellers}, nil
}
