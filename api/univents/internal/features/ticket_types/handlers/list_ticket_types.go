package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListTicketTypes(ctx context.Context, req openapi.ListTicketTypesRequestObject) (openapi.ListTicketTypesResponseObject, error) {
	ticketTypes, err := h.ops.ListByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListTicketTypes200JSONResponse{Items: ticketTypes}, nil
}
