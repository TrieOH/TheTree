package domain

import "context"

type TicketsRepository interface {
	Create(ctx context.Context, toCreate Ticket) (*Ticket, error)
}
