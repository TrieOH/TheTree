package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) ListWallets(ctx context.Context, _ openapi.ListWalletsRequestObject) (openapi.ListWalletsResponseObject, error) {
	wallets, err := h.ops.List(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListWallets201JSONResponse{Items: wallets}, nil
}
