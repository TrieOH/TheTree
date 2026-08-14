package checkouts

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
)

// GetEditionMyTicket is the "what do I hold" read: the caller's active
// ticket for the edition (pending or confirmed registration + its ticket
// type), or `data: null` when they hold none — the front uses it to show
// upgrade options on the more expensive ticket types. Read-only; the
// upgrade itself is a follow-up.
func (h *Handlers) GetEditionMyTicket(ctx context.Context, req openapi.GetEditionMyTicketRequestObject) (openapi.GetEditionMyTicketResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	res, err := h.ops.MyTicket(ctx, req.EditionId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	var data *openapi.MyTicket
	if res != nil {
		data = &openapi.MyTicket{
			RegistrationId: res.RegistrationID,
			Status:         openapi.MyTicketStatus(res.Status),
			TicketType:     res.TicketType,
		}
	}
	return openapi.GetEditionMyTicket200JSONResponse{
		Code:      200,
		Data:      data,
		Timestamp: time.Now(),
		Module:    module,
	}, nil
}
