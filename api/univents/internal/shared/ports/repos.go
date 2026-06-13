package ports

import (
	"context"
	"time"

	"univents/internal/shared/contracts"

	"github.com/google/uuid"
)

type EventsRepository interface {
	CreateEvent(ctx context.Context, toCreate *contracts.Event) (*contracts.Event, error)
	PatchEvent(ctx context.Context, toPatch *contracts.Event) (*contracts.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Event, error)
	ListEvents(ctx context.Context) ([]contracts.Event, error)
	ListOwnEvents(ctx context.Context, ownerID uuid.UUID) ([]contracts.Event, error)
	PublishEvent(ctx context.Context, id uuid.UUID) error
	AddEdition(ctx context.Context, eventID uuid.UUID) error
	AddGalleryImage(ctx context.Context, id uuid.UUID, url string) (*contracts.Event, error)
	RemoveGalleryImage(ctx context.Context, id uuid.UUID, url string) (*contracts.Event, error)
	SetLogo(ctx context.Context, id uuid.UUID, url string) (*contracts.Event, error)
	UnsetLogo(ctx context.Context, id uuid.UUID) (*contracts.Event, error)
	SetBanner(ctx context.Context, id uuid.UUID, url string) (*contracts.Event, error)
	UnsetBanner(ctx context.Context, id uuid.UUID) (*contracts.Event, error)
}

type EditionsRepository interface {
	Create(ctx context.Context, toCreate *contracts.Edition) (*contracts.Edition, error)
	GetByID(ctx context.Context, editionID uuid.UUID) (*contracts.Edition, error)
	List(ctx context.Context, editionID uuid.UUID) ([]contracts.Edition, error)
	ListAdmin(ctx context.Context, editionID uuid.UUID) ([]contracts.Edition, error)
	Announce(ctx context.Context, editionID uuid.UUID) error
	Open(ctx context.Context, editionID uuid.UUID) error
	Start(ctx context.Context, editionID uuid.UUID) error
	Finish(ctx context.Context, editionID uuid.UUID) error
	ConnectPaymentsAccount(ctx context.Context, editionID, triePaymentsCredentialID uuid.UUID, triePaymentsProvider, publicKey string) error
	DisconnectPaymentsAccount(ctx context.Context, editionID uuid.UUID) error
}

type ActivitiesRepository interface {
	Create(ctx context.Context, toCreate *contracts.Activity) (*contracts.Activity, error)
	Publish(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Activity, error)
	Start(ctx context.Context, id uuid.UUID) error
	Finish(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, editionID uuid.UUID) ([]contracts.Activity, error)
	ListAdmin(ctx context.Context, editionID uuid.UUID) ([]contracts.Activity, error)
	Register(ctx context.Context, toCreate contracts.AttendanceRecord) (*contracts.AttendanceRecord, error)
	Unregister(ctx context.Context, userID, activityID uuid.UUID) error
	MarkAttendanceRecordStatus(ctx context.Context, id uuid.UUID, scannedBy *uuid.UUID, status contracts.AttendanceStatus) error
	GetAttendanceRecordByID(ctx context.Context, id uuid.UUID) (*contracts.AttendanceRecord, error)
	ListActivityAttendanceRecords(ctx context.Context, activityID uuid.UUID) ([]contracts.AttendanceRecord, error)
	GetActiveUserActivityAttendanceRecords(ctx context.Context, userID, activityID uuid.UUID) (*contracts.AttendanceRecord, error)
	GetUserActivityAttendanceRecords(ctx context.Context, userID, activityID uuid.UUID) ([]contracts.AttendanceRecord, error)
	IsRegistered(ctx context.Context, userID, activityID uuid.UUID) (bool, error)
}

type CheckpointsRepository interface {
	Create(ctx context.Context, toCreate *contracts.Checkpoint) (*contracts.Checkpoint, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Checkpoint, error)
	List(ctx context.Context, editionID uuid.UUID) ([]contracts.Checkpoint, error)
}

type TicketsRepository interface {
	Create(ctx context.Context, toCreate contracts.Ticket) (*contracts.Ticket, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Ticket, error)
	AddPermission(ctx context.Context, toCreate contracts.TicketPermission) (*contracts.TicketPermission, error)
	RemovePermission(ctx context.Context, id, ticketID uuid.UUID) error
	List(ctx context.Context, editionID uuid.UUID) ([]contracts.Ticket, error)
	GetPermissions(ctx context.Context, ticketID uuid.UUID) ([]contracts.TicketPermission, error)
}

type ProductsRepository interface {
	Create(ctx context.Context, toCreate contracts.Product) (*contracts.Product, error)
	Publish(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Product, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]contracts.Product, error)
	List(ctx context.Context, editionID uuid.UUID) ([]contracts.Product, error)
	AdminList(ctx context.Context, editionID uuid.UUID) ([]contracts.Product, error)
	ReserveItems(ctx context.Context, sessionID uuid.UUID, items []contracts.CartItem, expiresAt time.Time) (contracts.ReservationOutcome, error)
	UnreserveItems(ctx context.Context, sessionID uuid.UUID) ([]contracts.InventoryUpdate, error)
	DeleteReservation(ctx context.Context, sessionID uuid.UUID) error
	Delete(ctx context.Context, productID uuid.UUID) error
	Restore(ctx context.Context, productID uuid.UUID) error
	ItemHasCompletedPurchases(ctx context.Context, productID uuid.UUID) (bool, error)
	AddGalleryImage(ctx context.Context, id uuid.UUID, url string) (*contracts.Product, error)
	RemoveGalleryImage(ctx context.Context, id uuid.UUID, url string) (*contracts.Product, error)
	SetThumbnail(ctx context.Context, id uuid.UUID, url string) (*contracts.Product, error)
	UnsetThumbnail(ctx context.Context, id uuid.UUID) (*contracts.Product, error)
}

type PurchaseRepository interface {
	Create(ctx context.Context, toCreate contracts.Purchase) (*contracts.Purchase, error)
	GetByPaymentID(ctx context.Context, paymentID string) (*contracts.Purchase, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*contracts.Purchase, error)
	CreateLineItem(ctx context.Context, toCreate contracts.LineItem) (*contracts.LineItem, error)
	ConfirmPurchase(ctx context.Context, paymentID string) error
	CancelPurchase(ctx context.Context, paymentID string) error
	GetTicketIDsByPaymentIntent(ctx context.Context, paymentID string) ([]contracts.TicketGrant, error)
	ListUserPurchases(ctx context.Context, userID uuid.UUID) ([]contracts.Purchase, error)
	ListPurchaseItems(ctx context.Context, purchaseID, userID uuid.UUID) ([]contracts.LineItem, error)
}

type InventoryPublisher interface {
	Publish(ctx context.Context, editionID uuid.UUID, updates []contracts.InventoryUpdate) error
}

type InventorySubscriber interface {
	Subscribe(ctx context.Context, editionID uuid.UUID) (<-chan []contracts.InventoryUpdate, error)
}

type PurchaseSessionStore interface {
	Save(ctx context.Context, session contracts.CheckoutSession) error
	Load(ctx context.Context, userID, sessionID uuid.UUID) (*contracts.CheckoutSession, error)
	Delete(ctx context.Context, userID, sessionID uuid.UUID) error
}
