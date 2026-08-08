package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type EventRepo interface {
	Create(ctx context.Context, toCreate *models.Event) (*models.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetBySlug(ctx context.Context, slug string) (*models.Event, error)
	ListPublic(ctx context.Context) ([]models.Event, error)
	ListOwned(ctx context.Context, ownerID uuid.UUID) ([]models.Event, error)
	ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Event, error)
	Publish(ctx context.Context, id uuid.UUID) error
	Discontinue(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, event *models.Event) (*models.Event, error)

	// Payments
	SetPaymentsConfig(ctx context.Context, id uuid.UUID, sellerID, walletID *uuid.UUID, publicKey *string) (*models.Event, error)
	ClearPaymentsConfig(ctx context.Context, id uuid.UUID) (*models.Event, error)

	// Members
	GetMember(ctx context.Context, eventID, userID uuid.UUID) (*models.EventMember, error)
	GetRole(ctx context.Context, actorID, eventID uuid.UUID) (models.EventMemberRole, error)
	AddEventMember(ctx context.Context, eventID, userID uuid.UUID, role models.EventMemberRole) (*models.EventMember, error)
	RemoveEventMember(ctx context.Context, eventID, userID uuid.UUID) error
	ListEventMembers(ctx context.Context, eventID uuid.UUID) ([]models.EventMember, error)
}
