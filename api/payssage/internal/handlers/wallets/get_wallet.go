package wallets

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) GetWallet(ctx context.Context, req openapi.GetWalletRequestObject) (openapi.GetWalletResponseObject, error) {
	wallet, err := h.ops.GetByID(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.GetWallet201JSONResponse{
		Code: 201, Data: wallet, Timestamp: time.Now(), Module: module,
	}, nil
}
