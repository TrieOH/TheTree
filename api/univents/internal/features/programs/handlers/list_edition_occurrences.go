package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionOccurrences(ctx context.Context, req openapi.ListEditionOccurrencesRequestObject) (openapi.ListEditionOccurrencesResponseObject, error) {
	occurrences, err := h.ops.ListOccurrencesByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionOccurrences200JSONResponse{Items: occurrences}, nil
}
