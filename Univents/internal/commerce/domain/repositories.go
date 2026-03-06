package domain

import (
	"context"

	"github.com/google/uuid"
)

type TicketsRepository interface {
	Create(ctx context.Context, toCreate Ticket) (*Ticket, error)
	AddPermission(ctx context.Context, toCreate TicketPermission) (*TicketPermission, error)
	RemovePermission(ctx context.Context, id, ticketID uuid.UUID) error
	List(ctx context.Context, editionID uuid.UUID) ([]Ticket, error)
}

type ProductsRepository interface {
	Create(ctx context.Context, toCreate Product) (*Product, error)
}
