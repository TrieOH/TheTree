package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListDraftEditions(ctx context.Context, req openapi.ListDraftEditionsRequestObject) (openapi.ListDraftEditionsResponseObject, error) {
	editions, err := h.ops.ListDraft(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListDraftEditions200JSONResponse{Items: editions}, nil
}
