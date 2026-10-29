package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListOwnedEvents(ctx context.Context, _ openapi.ListOwnedEventsRequestObject) (openapi.ListOwnedEventsResponseObject, error) {
	events, err := h.ops.ListOwned(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListOwnedEvents200JSONResponse{Items: events}, nil
}
