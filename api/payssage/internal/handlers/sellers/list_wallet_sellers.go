package sellers

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) ListWalletSellers(ctx context.Context, req openapi.ListWalletSellersRequestObject) (openapi.ListWalletSellersResponseObject, error) {
	sellers, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWalletSellers200JSONResponse{
		Code: 200, Data: &sellers, Timestamp: time.Now(), Module: module,
	}, nil
}
