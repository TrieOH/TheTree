package dto

import (
	"TriePayments/internal/core/domain"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateIntentRequest struct {
	Amount   int64           `json:"amount"`
	Currency string          `json:"currency"`
	Provider string          `json:"provider"`
	Metadata json.RawMessage `json:"metadata"`
}

type IntentResponse struct {
	ID           uuid.UUID       `json:"id"`
	WorkspaceID  uuid.UUID       `json:"workspace_id"`
	Amount       int64           `json:"amount"`
	Currency     string          `json:"currency"`
	Status       string          `json:"status"`
	ClientSecret string          `json:"client_secret"`
	Provider     string          `json:"provider"`
	Metadata     json.RawMessage `json:"metadata"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func MapIntentResponse(i *domain.Intent) IntentResponse {
	return IntentResponse{
		ID:           i.ID,
		WorkspaceID:  i.WorkspaceID,
		Amount:       i.Amount,
		Currency:     i.Currency,
		Status:       string(i.Status),
		ClientSecret: i.ClientSecret,
		Provider:     i.Provider,
		Metadata:     i.Metadata,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
}

type PayIntentRequest struct {
	CardToken       string `json:"card_token"       validate:"required"`
	PaymentMethodID string `json:"payment_method_id" validate:"required"`
	Installments    int    `json:"installments"      validate:"min=1"`
	PayerEmail      string `json:"payer_email"       validate:"required,email"`
}
