package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type EventRepo interface {
	Create(ctx context.Context, toCreate *models.Event) (*models.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	ListPublic(ctx context.Context) ([]models.Event, error)
	ListOwned(ctx context.Context, ownerID uuid.UUID) ([]models.Event, error)
	ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Event, error)
	Publish(ctx context.Context, id uuid.UUID) error
	Discontinue(ctx context.Context, id uuid.UUID) error

	// Members
	GetMember(ctx context.Context, eventID, userID uuid.UUID) (*models.EventMember, error)
	AddEventMember(ctx context.Context, eventID, userID uuid.UUID, role models.EventMemberRole) (*models.EventMember, error)
	RemoveEventMember(ctx context.Context, eventID, userID uuid.UUID) error
	ListEventMembers(ctx context.Context, eventID uuid.UUID) ([]models.EventMember, error)
}
