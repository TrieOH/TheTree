package handlers

import (
	"context"

	"payssage/internal/openapi"
)

// RefundIntent fully refunds a succeeded intent (webhook confirms the flip).
func (h *Handlers) RefundIntent(ctx context.Context, req openapi.RefundIntentRequestObject) (openapi.RefundIntentResponseObject, error) {
	intent, err := h.ops.Refund(ctx, req.IntentId)
	if err != nil {
		return nil, err
	}
	return openapi.RefundIntent200JSONResponse(*intent), nil
}
