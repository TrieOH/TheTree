package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) ListIntentsByProfile(ctx context.Context, _ openapi.ListIntentsByProfileRequestObject) (openapi.ListIntentsByProfileResponseObject, error) {
	intents, err := h.ops.ListByProfile(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListIntentsByProfile200JSONResponse{Items: intents}, nil
}
