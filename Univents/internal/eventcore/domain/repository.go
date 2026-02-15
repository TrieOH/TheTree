package domain

import (
	"context"

	"github.com/google/uuid"
)

type EventsRepository interface {
	Create(ctx context.Context, toCreate Event) (*Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	List(ctx context.Context) ([]Event, error)
	Publish(ctx context.Context, id uuid.UUID) error
	AppendAudit(ctx context.Context, audit Audit) (*Audit, error)
}
