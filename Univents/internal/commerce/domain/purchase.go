package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PurchaseStatus string

const (
	PurchaseStatusPending       PurchaseStatus = "pending"
	PurchaseStatusCompleted     PurchaseStatus = "completed"
	PurchaseStatusRefunded      PurchaseStatus = "refunded"
	PurchaseStatusPartialRefund PurchaseStatus = "partial_refund"
	PurchaseStatusCancelled     PurchaseStatus = "cancelled"
)

type Purchase struct {
	ID              uuid.UUID        `json:"id"`
	EditionID       uuid.UUID        `json:"edition_id"`
	UserID          uuid.UUID        `json:"user_id"`
	Status          PurchaseStatus   `json:"status"`
	SubtotalCents   int              `json:"subtotal_cents"`
	DiscountCents   int              `json:"discount_cents"`
	TaxCents        int              `json:"tax_cents"`
	TaxBreakdown    *json.RawMessage `json:"tax_breakdown"`
	TotalCents      int              `json:"total_cents"`
	PaymentProvider *string          `json:"payment_provider"`
	PaymentID       *string          `json:"payment_id"`
	FulfilledAt     *time.Time       `json:"fulfilled_at"`
	FulfilmentNotes *string          `json:"fulfilment_notes"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       *time.Time       `json:"deleted_at"`
}

type LineItem struct {
	ID                  uuid.UUID  `json:"id"`
	PurchaseID          uuid.UUID  `json:"purchase_id"`
	ItemType            string     `json:"item_type"`
	ItemID              uuid.UUID  `json:"item_id"`
	Quantity            int        `json:"quantity"`
	UnitPriceCents      int        `json:"unit_price_cents"`
	TotalPriceCents     int        `json:"total_price_cents"`
	AssignedToUserID    *uuid.UUID `json:"assigned_to_user_id"`
	Fulfilled           bool       `json:"fulfilled"`
	FulfilledAt         *time.Time `json:"fulfilled_at"`
	RefundedQuantity    int        `json:"refunded_quantity"`
	RefundedAmountCents int        `json:"refunded_amount_cents"`
	CreatedAt           time.Time  `json:"created_at"`
}

type CreatePurchaseSpec struct {
	EditionID       uuid.UUID        `json:"edition_id"`
	UserID          uuid.UUID        `json:"user_id"`
	SubtotalCents   int              `json:"subtotal_cents"`
	TaxCents        int              `json:"tax_cents"`
	TaxBreakdown    *json.RawMessage `json:"tax_breakdown"`
	TotalCents      int              `json:"total_cents"`
	PaymentProvider *string          `json:"payment_provider"`
	PaymentID       *string          `json:"payment_id"`
}

func NewPurchase(spec CreatePurchaseSpec) *Purchase {
	return &Purchase{
		ID:              uuid.UUID{},
		EditionID:       spec.EditionID,
		UserID:          spec.UserID,
		Status:          PurchaseStatusPending,
		SubtotalCents:   spec.SubtotalCents,
		DiscountCents:   0,
		TaxCents:        spec.TaxCents,
		TaxBreakdown:    spec.TaxBreakdown,
		TotalCents:      spec.TotalCents,
		PaymentProvider: spec.PaymentProvider,
		PaymentID:       spec.PaymentID,
	}
}
