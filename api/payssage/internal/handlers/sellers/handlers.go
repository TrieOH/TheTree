// Package sellers implements the StrictServerInterface methods for the
// sellers feature.
package sellers

import (
	"context"
	"time"

	"payssage/internal/openapi"
	"payssage/internal/services"
)

const module = "Payssage"

type Handlers struct {
	ops *services.Sellers
}

func New(ops *services.Sellers) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListWalletSellers(ctx context.Context, req openapi.ListWalletSellersRequestObject) (openapi.ListWalletSellersResponseObject, error) {
	sellers, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWalletSellers200JSONResponse{
		Code: 200, Data: &sellers, Timestamp: time.Now(), Module: module,
	}, nil
}
