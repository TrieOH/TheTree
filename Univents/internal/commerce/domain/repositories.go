package domain

import "context"

type TicketsRepository interface {
	Create(ctx context.Context, toCreate Ticket) (*Ticket, error)
}

type ProductsRepository interface {
	Create(ctx context.Context, toCreate Product) (*Product, error)
}
