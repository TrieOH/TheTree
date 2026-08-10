package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID                   uuid.UUID  `json:"id"`
	EditionID            uuid.UUID  `json:"edition_id"`
	VendorCode           string     `json:"vendor_code"`
	RequiresRegistration bool       `json:"requires_registration"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at"`
}

// ProductPurchase is one buyer's purchase of a product variant at an
// edition (a purchase_items materialization, D4). Created pending at
// checkout; flipped confirmed/cancelled/expired by the webhook receiver
// (split 4) and the expiry worker (split 7).
type ProductPurchaseStatus string

const (
	ProductPurchaseStatusPending   ProductPurchaseStatus = "pending"
	ProductPurchaseStatusConfirmed ProductPurchaseStatus = "confirmed"
	ProductPurchaseStatusCancelled ProductPurchaseStatus = "cancelled"
	ProductPurchaseStatusExpired   ProductPurchaseStatus = "expired"
)

type ProductPurchase struct {
	ID               uuid.UUID             `json:"id"`
	EditionID        uuid.UUID             `json:"edition_id"`
	VariantID        uuid.UUID             `json:"variant_id"`
	PurchaserID      uuid.UUID             `json:"purchaser_id"`
	RecipientID      *uuid.UUID            `json:"recipient_id"`
	RegistrationID   *uuid.UUID            `json:"registration_id"`
	Quantity         int                   `json:"quantity"`
	Status           ProductPurchaseStatus `json:"status"`
	StatusReason     *string               `json:"status_reason"`
	PayssageIntentID *uuid.UUID            `json:"payssage_intent_id"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        *time.Time            `json:"updated_at"`
	DeletedAt        *time.Time            `json:"deleted_at"`
}

type ProductVariant struct {
	ID          uuid.UUID       `json:"id"`
	EditionID   uuid.UUID       `json:"edition_id"`
	ProductID   uuid.UUID       `json:"product_id"`
	VendorCode  string          `json:"vendor_code"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Price       int64           `json:"price"`
	Stock       *int            `json:"stock"`
	GalleryURLs json.RawMessage `json:"gallery_urls"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at"`
	DeletedAt   *time.Time      `json:"deleted_at"`
}

// ── Create Initial Product (product + first variant) ──────────────────────

type CreateInitialProductInput struct {
	EditionID            uuid.UUID `json:"edition_id"`
	VendorCode           string    `json:"vendor_code"`
	VariantVendorCode    string    `json:"variant_vendor_code"`
	RequiresRegistration bool      `json:"requires_registration"`
	Name                 string    `json:"name"`
	Description          *string   `json:"description"`
	Price                int64     `json:"price"`
	Stock                *int      `json:"stock"`
}

// ── Create Product Variant ────────────────────────────────────────────────

type CreateProductVariantInput struct {
	ProductID   uuid.UUID `json:"product_id"`
	VendorCode  string    `json:"vendor_code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       int64     `json:"price"`
	Stock       *int      `json:"stock"`
}

// ── Patch Product ─────────────────────────────────────────────────────────

type PatchProductInput struct {
	ProductID            uuid.UUID
	VendorCode           string
	RequiresRegistration bool
}

// ── Patch Product Variant ─────────────────────────────────────────────────

type PatchProductVariantInput struct {
	VariantID   uuid.UUID
	VendorCode  string
	Name        string
	Description *string
	Price       int64
	Stock       *int
	GalleryURLs json.RawMessage
}
