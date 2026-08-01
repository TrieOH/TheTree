package wallets

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) UnbindCollector(ctx context.Context, req openapi.UnbindCollectorRequestObject) (openapi.UnbindCollectorResponseObject, error) {
	err := h.ops.UnbindCollector(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.UnbindCollector204Response{}, nil
}
