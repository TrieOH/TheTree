package wallets

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) BindCollector(ctx context.Context, req openapi.BindCollectorRequestObject) (openapi.BindCollectorResponseObject, error) {
	err := h.ops.BindCollector(ctx, req.WalletId, req.Body.CollectorId)
	if err != nil {
		return nil, err
	}
	return openapi.BindCollector204Response{}, nil
}
