package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseStatus string

const (
	PurchaseStatusPending   PurchaseStatus = "pending"
	PurchaseStatusApproved  PurchaseStatus = "approved"
	PurchaseStatusExpired   PurchaseStatus = "expired"
	PurchaseStatusCancelled PurchaseStatus = "cancelled"
)

type PurchaseItemType string

const (
	PurchaseItemTypeTicket            PurchaseItemType = "ticket"
	PurchaseItemTypeProduct           PurchaseItemType = "product"
	PurchaseItemTypeProgramOccurrence PurchaseItemType = "program_occurrence"
)

// Purchase is one Edition-scoped order (tickets/products/program spots),
// created at checkout. `PayssageIntentID` is the correlation key — the
// webhook receiver matches on it, never on a provider-specific id (D2).
// `purchase_items` is the availability ledger; materialized rows
// (registrations/product_purchases/program_participations) follow it in the
// same tx (D4) and are flipped on approve/cancel/expire.
type Purchase struct {
	ID               uuid.UUID      `json:"id"`
	EditionID        uuid.UUID      `json:"edition_id"`
	PurchaserID      uuid.UUID      `json:"purchaser_id"`
	Status           PurchaseStatus `json:"status"`
	StatusReason     *string        `json:"status_reason"`
	TotalCents       int64          `json:"total_cents"`
	Currency         string         `json:"currency"`
	PaymentMethod    *string        `json:"payment_method"`
	PayssageSellerID *uuid.UUID     `json:"payssage_seller_id"`
	PayssageIntentID *uuid.UUID     `json:"payssage_intent_id"`
	QRCode           *string        `json:"qr_code"`
	QRCodeBase64     *string        `json:"qr_code_base64"`
	ExpiresAt        time.Time      `json:"expires_at"`
	RiverJobID       *int64         `json:"river_job_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        *time.Time     `json:"updated_at"`
	DeletedAt        *time.Time     `json:"deleted_at"`
}

type PurchaseItem struct {
	ID                uuid.UUID        `json:"id"`
	PurchaseID        uuid.UUID        `json:"purchase_id"`
	ItemType          PurchaseItemType `json:"item_type"`
	ItemID            uuid.UUID        `json:"item_id"`
	Quantity          int              `json:"quantity"`
	UnitPriceCents    int64            `json:"unit_price_cents"`
	RegistrationID    *uuid.UUID       `json:"registration_id"`
	ProductPurchaseID *uuid.UUID       `json:"product_purchase_id"`
	ParticipationID   *uuid.UUID       `json:"participation_id"`
}

// WsToken is a one-time handshake token for the per-purchase WebSocket
// (split 6): it proves prior REST auth for this purchase and nothing more.
// Stored as a SHA-256 hash; short-lived (10min) and single-use.
type WsToken struct {
	ID         uuid.UUID  `json:"id"`
	PurchaseID uuid.UUID  `json:"purchase_id"`
	UserID     uuid.UUID  `json:"user_id"`
	TokenHash  string     `json:"token_hash"`
	ExpiresAt  time.Time  `json:"expires_at"`
	UsedAt     *time.Time `json:"used_at"`
}

// ItemAvailability is one purchasable item's stock position:
// available = base - reserved; nil base = unlimited.
type ItemAvailability struct {
	ItemType         PurchaseItemType `json:"item_type"`
	ItemID           uuid.UUID        `json:"item_id"`
	BaseQuantity     *int             `json:"base_quantity"`
	ReservedQuantity int64            `json:"reserved_quantity"`
}
