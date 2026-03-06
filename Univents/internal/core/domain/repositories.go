package domain

import (
	"context"

	"github.com/google/uuid"
)

type EventsRepository interface {
	CreateEvent(ctx context.Context, toCreate *Event) (*Event, error)
	PatchEvent(ctx context.Context, toPatch *Event) (*Event, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*Event, error)
	ListEvents(ctx context.Context) ([]Event, error)
	ListOwnEvents(ctx context.Context, ownerID uuid.UUID) ([]Event, error)
	PublishEvent(ctx context.Context, id uuid.UUID) error
	AddEdition(ctx context.Context, eventID uuid.UUID) error
}

type EditionsRepository interface {
	Create(ctx context.Context, toCreate *Edition) (*Edition, error)
	GetByID(ctx context.Context, editionID uuid.UUID) (*Edition, error)
	List(ctx context.Context, editionID uuid.UUID) ([]Edition, error)
	Announce(ctx context.Context, editionID uuid.UUID) error
	Open(ctx context.Context, editionID uuid.UUID) error
	Start(ctx context.Context, editionID uuid.UUID) error
	Finish(ctx context.Context, editionID uuid.UUID) error
}

type ActivitiesRepository interface {
	Create(ctx context.Context, toCreate *Activity) (*Activity, error)
	Publish(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Activity, error)
	Start(ctx context.Context, id uuid.UUID) error
	Finish(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, editionID uuid.UUID) ([]Activity, error)
	ListAdmin(ctx context.Context, editionID uuid.UUID) ([]Activity, error)
}

type CheckpointsRepository interface {
	Create(ctx context.Context, toCreate *Checkpoint) (*Checkpoint, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Checkpoint, error)
	List(ctx context.Context, editionID uuid.UUID) ([]Checkpoint, error)
}
