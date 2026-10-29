package ticket_types

import (
	"context"

	"univents/internal/features/ticket_types/service"
	"univents/internal/openapi"
	"univents/models"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

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
	return openapi.CreateTicketType201JSONResponse(*ticketType), nil
}

func (h *Handlers) GetTicketType(ctx context.Context, req openapi.GetTicketTypeRequestObject) (openapi.GetTicketTypeResponseObject, error) {
	ticketType, err := h.ops.GetByID(ctx, req.TicketTypeId)
	if err != nil {
		return nil, err
	}
	return openapi.GetTicketType200JSONResponse(*ticketType), nil
}

func (h *Handlers) ListTicketTypes(ctx context.Context, req openapi.ListTicketTypesRequestObject) (openapi.ListTicketTypesResponseObject, error) {
	ticketTypes, err := h.ops.ListByEdition(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListTicketTypes200JSONResponse{Items: ticketTypes}, nil
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
	return openapi.PatchTicketType200JSONResponse(*ticketType), nil
}
