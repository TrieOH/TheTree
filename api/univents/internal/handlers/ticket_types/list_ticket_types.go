package ticket_types

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListTicketTypes(ctx context.Context, req openapi.ListTicketTypesRequestObject) (openapi.ListTicketTypesResponseObject, error) {
	ticketTypes, err := h.ops.ListByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListTicketTypes200JSONResponse{
		Code: 200, Data: &ticketTypes, Timestamp: time.Now(), Module: module,
	}, nil
}
