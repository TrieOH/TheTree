package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateIntentRequest struct {
	Amount             int64           `json:"amount"`
	Currency           string          `json:"currency"`
	Provider           string          `json:"provider"`
	Metadata           json.RawMessage `json:"metadata"`
	PaymentMethodID    string          `json:"payment_method_id"`
	Installments       int             `json:"installments"`
	CardToken          string          `json:"card_token"`
	PaymentMethodType  string          `json:"payment_method_type"`
	SellerCredentialID uuid.UUID       `json:"seller_credential_id"`
	PayerEmail         string          `json:"payer_email"`
}

type PayIntentRequest struct {
	SellerCredentialID uuid.UUID `json:"seller_credential_id"`
}
