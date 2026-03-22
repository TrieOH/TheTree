package domain

import (
	"TriePayments/internal/shared/errx"
	"TriePayments/internal/shared/validation"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Intent struct {
	ID                 uuid.UUID       `json:"id"`
	WorkspaceID        uuid.UUID       `json:"workspace_id"`
	Amount             int64           `json:"amount"`
	Currency           string          `json:"currency"`
	Status             IntentStatus    `json:"status"`
	Provider           string          `json:"provider"`
	Metadata           json.RawMessage `json:"metadata"`
	SellerCredentialID *uuid.UUID      `json:"seller_credential_id,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`

	// Only one of these will be non-nil, determined by Provider.
	MercadoPagoData *MercadoPagoIntentData `json:"mercadopago_data,omitempty"`
}

type MercadoPagoIntentData struct {
	OrderID                 string `json:"order_id"`
	OrderStatus             string `json:"order_status"`
	OrderStatusDetail       string `json:"order_status_detail"`
	TransactionID           string `json:"transaction_id"`
	TransactionStatus       string `json:"transaction_status"`
	TransactionStatusDetail string `json:"transaction_status_detail"`
	PaymentMethodID         string `json:"payment_method_id"`
	PaymentMethodType       string `json:"payment_method_type"`
	PixQRCode               string `json:"pix_qr_code,omitempty"` //FIXME maybe dont send this or the one below
	PixQRCodeB64            string `json:"pix_qr_code_base64,omitempty"`
}

type IntentStatus string

const (
	IntentStatusPending   IntentStatus = "pending"
	IntentStatusSucceeded IntentStatus = "succeeded"
	IntentStatusCancelled IntentStatus = "cancelled"
	IntentStatusFailed    IntentStatus = "failed"
)

func NewIntent(workspaceID uuid.UUID, amount int64, currency, provider string, metadata json.RawMessage) (*Intent, error) {
	if metadata == nil {
		metadata = json.RawMessage("{}")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("product").SetMessage("error generating uuid").SetCause(err)
	}

	i := &Intent{
		ID:          id,
		WorkspaceID: workspaceID,
		Amount:      amount,
		Currency:    currency,
		Status:      IntentStatusPending,
		Provider:    provider,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := i.validate(); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Intent) validate() error {
	return validation.Run(
		validation.RequireUUID("intent", "workspace_id", i.WorkspaceID),
		validation.RequireString("intent", "currency", i.Currency),
		validation.RequireString("intent", "provider", i.Provider),
		validation.Assert("intent", i.Amount > 0, "amount must be greater than zero"),
		validation.Assert("intent", len(i.Currency) == 3, "currency must be a 3-letter ISO code"),
	)
}
