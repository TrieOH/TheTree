package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListUpcomingEditions(ctx context.Context, req openapi.ListUpcomingEditionsRequestObject) (openapi.ListUpcomingEditionsResponseObject, error) {
	editions, err := h.ops.GetUpcoming(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListUpcomingEditions200JSONResponse{Items: editions}, nil
}
