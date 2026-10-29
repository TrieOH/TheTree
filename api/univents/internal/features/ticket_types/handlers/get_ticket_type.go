package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) GetTicketType(ctx context.Context, req openapi.GetTicketTypeRequestObject) (openapi.GetTicketTypeResponseObject, error) {
	ticketType, err := h.ops.GetByID(ctx, req.TicketTypeId)
	if err != nil {
		return nil, err
	}
	return openapi.GetTicketType200JSONResponse(*ticketType), nil
}
