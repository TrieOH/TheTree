package ticket_types

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) GetTicketType(ctx context.Context, req openapi.GetTicketTypeRequestObject) (openapi.GetTicketTypeResponseObject, error) {
	ticketType, err := h.ops.GetByID(ctx, req.TicketTypeId)
	if err != nil {
		return nil, err
	}
	return openapi.GetTicketType200JSONResponse{
		Code: 200, Data: ticketType, Timestamp: time.Now(), Module: module,
	}, nil
}
