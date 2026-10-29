package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListEditionBadgeEmissions(ctx context.Context, req openapi.ListEditionBadgeEmissionsRequestObject) (openapi.ListEditionBadgeEmissionsResponseObject, error) {
	items, err := h.ops.ListByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEditionBadgeEmissions200JSONResponse{Items: items}, nil
}
