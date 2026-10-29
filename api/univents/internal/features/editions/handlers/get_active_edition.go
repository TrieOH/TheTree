package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) GetActiveEdition(ctx context.Context, req openapi.GetActiveEditionRequestObject) (openapi.GetActiveEditionResponseObject, error) {
	edition, err := h.ops.GetActive(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActiveEdition200JSONResponse(*edition), nil
}
