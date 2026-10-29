package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListPastEditions(ctx context.Context, req openapi.ListPastEditionsRequestObject) (openapi.ListPastEditionsResponseObject, error) {
	editions, err := h.ops.GetPast(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListPastEditions200JSONResponse{Items: editions}, nil
}
