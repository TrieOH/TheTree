package dtos

import (
	"time"
	"univents/internal/commerce/domain"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	EditionScopeID    uuid.UUID          `json:"edition_scope_id"`
	Name              string             `json:"name" validate:"required,min=3"`
	Description       *string            `json:"description"`
	Type              domain.ProductType `json:"type"`
	TicketID          *uuid.UUID         `json:"ticket_id"`
	PriceCents        int                `json:"price_cents" validate:"gte=0"`
	AvailableFrom     *time.Time         `json:"available_from"`
	AvailableUntil    *time.Time         `json:"available_until"`
	HasInventory      bool               `json:"has_inventory"`
	InventoryQuantity int                `json:"inventory_quantity"`
}

type BuyRequest struct {
	Items []domain.CartItem `json:"items"`
}

type ReservationConfirmedPayload struct {
	SessionID       uuid.UUID `json:"session_id"`
	ClientSecret    string    `json:"client_secret"`
	PaymentIntentID string    `json:"payment_intent_id"`
	ExpiresAt       time.Time `json:"expires_at"`
}

type ConfirmPaymentRequest struct {
	SessionID       uuid.UUID `json:"session_id"`
	PaymentIntentID string    `json:"payment_intent_id"`
}
