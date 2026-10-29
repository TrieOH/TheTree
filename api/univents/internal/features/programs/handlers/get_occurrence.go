package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) GetOccurrence(ctx context.Context, req openapi.GetOccurrenceRequestObject) (openapi.GetOccurrenceResponseObject, error) {
	occurrence, err := h.ops.GetOccurrenceByID(ctx, req.OccurrenceId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOccurrence200JSONResponse(*occurrence), nil
}
