package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) DeleteOccurrence(ctx context.Context, req openapi.DeleteOccurrenceRequestObject) (openapi.DeleteOccurrenceResponseObject, error) {
	occurrence, err := h.ops.DeleteOccurrence(ctx, req.OccurrenceId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteOccurrence200JSONResponse(*occurrence), nil
}
