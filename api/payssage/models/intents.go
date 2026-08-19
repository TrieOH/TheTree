package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Intent struct {
	ID           uuid.UUID           `json:"id"`
	WalletID     uuid.UUID           `json:"wallet_id"`
	SellerID     uuid.UUID           `json:"seller_id"`
	CollectorID  *uuid.UUID          `json:"collector_id"`
	AmountCents  int64               `json:"amount_cents"`
	Currency     string              `json:"currency"`
	Sandbox      bool                `json:"sandbox"`
	Provider     string              `json:"provider"`
	Status       IntentStatus        `json:"status"`
	StatusDetail *IntentStatusDetail `json:"status_detail"`
	ProviderData json.RawMessage     `json:"provider_data"`
	Metadata     *json.RawMessage    `json:"metadata"`
	// ExternalID is the caller's order id (e.g. the univents purchase id).
	// ExternalGroup groups intents by a caller-side grouping key (e.g. the
	// univents edition id). Both are counterparty ids — see the naming
	// invariant in the refund plan (A7): external_* always holds the
	// counterparty's id, never payssage's own.
	ExternalID    *string   `json:"external_id,omitempty"`
	ExternalGroup *string   `json:"external_group,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type IntentStatus string

const (
	IntentStatusPending    IntentStatus = "pending"
	IntentStatusProcessing IntentStatus = "processing"
	IntentStatusSucceeded  IntentStatus = "succeeded"
	IntentStatusCancelled  IntentStatus = "cancelled"
	IntentStatusRejected   IntentStatus = "rejected"
	IntentStatusFailed     IntentStatus = "failed"
	IntentStatusRefunded   IntentStatus = "refunded"
)

type IntentStatusDetail string

const (
	StatusDetailInsufficientFunds   IntentStatusDetail = "insufficient_funds"
	StatusDetailHighRisk            IntentStatusDetail = "high_risk"
	StatusDetailInvalidCard         IntentStatusDetail = "invalid_card"
	StatusDetailCardDisabled        IntentStatusDetail = "card_disabled"
	StatusDetailExpiredCard         IntentStatusDetail = "expired_card"
	StatusDetailInvalidSecurityCode IntentStatusDetail = "invalid_security_code"
	StatusDetailPendingReview       IntentStatusDetail = "pending_review"
	StatusDetailOther               IntentStatusDetail = "other"
)

type CreateIntentInput struct {
	WalletID      uuid.UUID
	SellerID      uuid.UUID
	Currency      string
	AmountCents   int64
	CheckoutData  json.RawMessage
	Metadata      *json.RawMessage
	ExternalID    *string
	ExternalGroup *string
}

type HardCreateIntentInput struct {
	WalletID      uuid.UUID        `json:"wallet_id"`
	SellerID      uuid.UUID        `json:"seller_id"`
	CollectorID   *uuid.UUID       `json:"collector_id"`
	AmountCents   int64            `json:"amount_cents"`
	Currency      string           `json:"currency"`
	Sandbox       bool             `json:"sandbox"`
	Provider      string           `json:"provider"`
	Status        IntentStatus     `json:"status"`
	ProviderData  json.RawMessage  `json:"provider_data"`
	Metadata      *json.RawMessage `json:"metadata"`
	ExternalID    *string          `json:"external_id"`
	ExternalGroup *string          `json:"external_group"`
}
