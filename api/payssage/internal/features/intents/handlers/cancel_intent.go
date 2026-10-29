package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) CancelIntent(ctx context.Context, req openapi.CancelIntentRequestObject) (openapi.CancelIntentResponseObject, error) {
	intent, err := h.ops.Cancel(ctx, req.IntentId)
	if err != nil {
		return nil, err
	}
	return openapi.CancelIntent200JSONResponse(*intent), nil
}
