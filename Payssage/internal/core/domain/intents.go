package domain

import (
	"TriePayments/internal/shared/errx"
	"TriePayments/internal/shared/validation"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Intent struct {
	ID                uuid.UUID       `json:"id"`
	WorkspaceID       uuid.UUID       `json:"workspace_id"`
	Amount            int64           `json:"amount"`
	Currency          string          `json:"currency"`
	Status            IntentStatus    `json:"status"`
	ClientSecret      string          `json:"client_secret"`
	Provider          string          `json:"provider"`
	ProviderPaymentID *string         `json:"provider_payment_id"`
	Metadata          json.RawMessage `json:"metadata"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
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
		ID:           id,
		WorkspaceID:  workspaceID,
		Amount:       amount,
		Currency:     currency,
		Status:       IntentStatusPending,
		ClientSecret: "secret_mock_" + uuid.NewString(),
		Provider:     provider,
		Metadata:     metadata,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
