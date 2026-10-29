package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) GetWallet(ctx context.Context, req openapi.GetWalletRequestObject) (openapi.GetWalletResponseObject, error) {
	wallet, err := h.ops.GetByID(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWallet201JSONResponse(*wallet), nil
}
