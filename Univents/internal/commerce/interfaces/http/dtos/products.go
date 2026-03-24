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

type ResumeSessionPayload struct {
	SessionID uuid.UUID `json:"session_id"`
}

type BuyRequest struct {
	Items []domain.CartItem `json:"items"`
}

type ReservationConfirmedPayload struct {
	SessionID     uuid.UUID             `json:"session_id"`
	ExpiresAt     time.Time             `json:"expires_at"`
	ReservedItems []domain.ReservedItem `json:"reserved_items"`
	TotalCents    int                   `json:"total_cents"`
}

type ConfirmPaymentRequest struct {
	SessionID       uuid.UUID `json:"session_id"`
	PaymentIntentID string    `json:"payment_intent_id"`
}

type SubmitPaymentPayload struct {
	CardToken            string `json:"card_token"`
	PaymentMethodID      string `json:"payment_method_id"`
	PaymentMethodType    string `json:"payment_method_type"`
	Installments         int    `json:"installments"`
	PayerEmail           string `json:"payer_email"`
	IdentificationNumber string `json:"identification_number"`
	IdentificationType   string `json:"identification_type"`
}

type ImageURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}
