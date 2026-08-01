package ticket_types

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) CreateTicketType(ctx context.Context, req openapi.CreateTicketTypeRequestObject) (openapi.CreateTicketTypeResponseObject, error) {
	ticketType, err := h.ops.Create(ctx, models.CreateTicketTypeInput{
		EditionID:   req.EditionId,
		Name:        req.Body.Name,
		Description: req.Body.Description,
		AccessLevel: req.Body.AccessLevel,
		PriceCents:  req.Body.PriceCents,
		MaxQuantity: req.Body.MaxQuantity,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateTicketType201JSONResponse{
		Code: 201, Data: ticketType, Timestamp: time.Now(), Module: module,
	}, nil
}
