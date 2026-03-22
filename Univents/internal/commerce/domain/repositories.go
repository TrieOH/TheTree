package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TicketsRepository interface {
	Create(ctx context.Context, toCreate Ticket) (*Ticket, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Ticket, error)
	AddPermission(ctx context.Context, toCreate TicketPermission) (*TicketPermission, error)
	RemovePermission(ctx context.Context, id, ticketID uuid.UUID) error
	List(ctx context.Context, editionID uuid.UUID) ([]Ticket, error)
	GetPermissions(ctx context.Context, ticketID uuid.UUID) ([]TicketPermission, error)
}

type ProductsRepository interface {
	Create(ctx context.Context, toCreate Product) (*Product, error)
	Publish(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]Product, error)
	List(ctx context.Context, editionID uuid.UUID) ([]Product, error)
	AdminList(ctx context.Context, editionID uuid.UUID) ([]Product, error)
	ReserveItems(ctx context.Context, sessionID uuid.UUID, items []CartItem, expiresAt time.Time) (ReservationOutcome, error)
	UnreserveItems(ctx context.Context, sessionID uuid.UUID) ([]InventoryUpdate, error)
	DeleteReservation(ctx context.Context, sessionID uuid.UUID) error
}

type PurchaseRepository interface {
	Create(ctx context.Context, toCreate Purchase) (*Purchase, error)
	GetByPaymentID(ctx context.Context, paymentID string) (*Purchase, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*Purchase, error)
	CreateLineItem(ctx context.Context, toCreate LineItem) (*LineItem, error)
	ConfirmPurchase(ctx context.Context, paymentID string) error
	CancelPurchase(ctx context.Context, paymentID string) error
	GetTicketIDsByPaymentIntent(ctx context.Context, paymentID string) ([]TicketGrant, error)
	ListUserPurchases(ctx context.Context, userID uuid.UUID) ([]Purchase, error)
	ListPurchaseItems(ctx context.Context, purchaseID, userID uuid.UUID) ([]LineItem, error)
}

type InventoryPublisher interface {
	Publish(ctx context.Context, editionID uuid.UUID, updates []InventoryUpdate) error
}

type InventorySubscriber interface {
	Subscribe(ctx context.Context, editionID uuid.UUID) (<-chan []InventoryUpdate, error)
}

type PurchaseSessionStore interface {
	Save(ctx context.Context, session PurchaseSession) error
	Load(ctx context.Context, userID, sessionID uuid.UUID) (*PurchaseSession, error)
	Delete(ctx context.Context, userID, sessionID uuid.UUID) error
}
