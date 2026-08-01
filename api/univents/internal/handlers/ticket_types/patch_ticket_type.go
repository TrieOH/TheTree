package ticket_types

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) PatchTicketType(ctx context.Context, req openapi.PatchTicketTypeRequestObject) (openapi.PatchTicketTypeResponseObject, error) {
	ticketType, err := h.ops.Patch(ctx, models.PatchTicketTypeInput{
		TicketTypeID: req.TicketTypeId,
		Name:         req.Body.Name,
		Description:  req.Body.Description,
		AccessLevel:  req.Body.AccessLevel,
		PriceCents:   req.Body.PriceCents,
		MaxQuantity:  req.Body.MaxQuantity,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PatchTicketType200JSONResponse{
		Code: 200, Data: ticketType, Timestamp: time.Now(), Module: module,
	}, nil
}
