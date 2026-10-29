package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListPublicEvents(ctx context.Context, _ openapi.ListPublicEventsRequestObject) (openapi.ListPublicEventsResponseObject, error) {
	events, err := h.ops.ListPublic(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListPublicEvents200JSONResponse{Items: events}, nil
}
