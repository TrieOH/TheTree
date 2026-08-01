package intents

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) GetIntent(ctx context.Context, req openapi.GetIntentRequestObject) (openapi.GetIntentResponseObject, error) {
	intent, err := h.ops.GetByID(ctx, req.IntentId)
	if err != nil {
		return nil, err
	}
	return openapi.GetIntent200JSONResponse{
		Code: 200, Data: intent, Timestamp: time.Now(), Module: module,
	}, nil
}
