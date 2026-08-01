package intents

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) ListIntentsByProfile(ctx context.Context, _ openapi.ListIntentsByProfileRequestObject) (openapi.ListIntentsByProfileResponseObject, error) {
	intents, err := h.ops.ListByProfile(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListIntentsByProfile200JSONResponse{
		Code: 200, Data: &intents, Timestamp: time.Now(), Module: module,
	}, nil
}
