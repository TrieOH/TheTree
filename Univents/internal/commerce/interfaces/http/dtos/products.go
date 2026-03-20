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
	SessionID uuid.UUID        `json:"session_id"`
	ExpiresAt time.Time        `json:"expires_at"`
	Items     []map[string]any `json:"items"`
	Total     int              `json:"total"`
}

type ConfirmPaymentRequest struct {
	SessionID       uuid.UUID `json:"session_id"`
	PaymentIntentID string    `json:"payment_intent_id"`
}

type SubmitPaymentPayload struct {
	CardToken          string `json:"card_token"`
	PaymentMethodID    string `json:"payment_method_id"`
	PaymentMethodType  string `json:"payment_method_type"`
	Installments       int    `json:"installments"`
	SellerCredentialID string
	PayerEmail         string
}

type OrderPayload struct {
	CardToken       string `json:"card_token"`
	PaymentMethodID string `json:"payment_method_id"`
	Installments    int    `json:"installments"`
}
