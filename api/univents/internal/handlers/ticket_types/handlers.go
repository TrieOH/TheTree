// Package ticket_types implements the StrictServerInterface methods for
// the ticket_types feature.
package ticket_types

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/internal/services"
	"univents/models"
)

const module = "Univents"

type Handlers struct {
	ops *services.TicketTypes
}

func New(ops *services.TicketTypes) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListTicketTypes(ctx context.Context, req openapi.ListTicketTypesRequestObject) (openapi.ListTicketTypesResponseObject, error) {
	ticketTypes, err := h.ops.ListByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListTicketTypes200JSONResponse{
		Code: 200, Data: &ticketTypes, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetTicketType(ctx context.Context, req openapi.GetTicketTypeRequestObject) (openapi.GetTicketTypeResponseObject, error) {
	ticketType, err := h.ops.GetByID(ctx, req.TicketTypeId)
	if err != nil {
		return nil, err
	}
	return openapi.GetTicketType200JSONResponse{
		Code: 200, Data: ticketType, Timestamp: time.Now(), Module: module,
	}, nil
}

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
