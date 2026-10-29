package sellers

import (
	"context"

	"payssage/internal/features/sellers/service"
	"payssage/internal/openapi"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListWalletSellers(ctx context.Context, req openapi.ListWalletSellersRequestObject) (openapi.ListWalletSellersResponseObject, error) {
	sellers, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWalletSellers200JSONResponse{Items: sellers}, nil
}
