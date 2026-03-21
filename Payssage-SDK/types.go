package paymentsSDK

import (
	"encoding/json"
	"time"
)

type ProviderCredential struct {
	ID          string     `json:"id"`
	WorkspaceID string     `json:"workspace_id"`
	Provider    string     `json:"provider"`
	DisplayName string     `json:"display_name"`
	CreatedAt   time.Time  `json:"created_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
}

type Workspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type APIKey struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Prefix    string     `json:"prefix"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

type CreateAPIKeyResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Key       string     `json:"key"` // only returned once
	Prefix    string     `json:"prefix"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

type IntentStatus string

const (
	IntentStatusPending   IntentStatus = "pending"
	IntentStatusSucceeded IntentStatus = "succeeded"
	IntentStatusCancelled IntentStatus = "cancelled"
	IntentStatusFailed    IntentStatus = "failed"
)

type Intent struct {
	ID          string          `json:"id"`
	WorkspaceID string          `json:"workspace_id"`
	Amount      int64           `json:"amount"`
	Currency    string          `json:"currency"`
	Status      IntentStatus    `json:"status"`
	Provider    string          `json:"provider"`
	Metadata    json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`

	// Only one of these will be non-nil, determined by Provider.
	MercadoPagoData *MercadoPagoIntentData `json:"mercadopago_data,omitempty"`
}

type MercadoPagoIntentData struct {
	OrderID           string  `json:"order_id"`
	OrderStatus       string  `json:"order_status"`
	OrderStatusDetail string  `json:"order_status_detail"`
	TransactionID     string  `json:"transaction_id"`
	PixQRCode         *string `json:"pix_qr_code,omitempty"` //FIXME maybe dont send this or the one below
	PixQRCodeB64      *string `json:"pix_qr_code_base64,omitempty"`
}

type WebhookEndpoint struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateWebhookEndpointResponse struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	URL         string    `json:"url"`
	Secret      string    `json:"secret"` // only returned on creation
	CreatedAt   time.Time `json:"created_at"`
}

// WebhookPayload is the normalized event delivered to registered endpoints
type WebhookPayload struct {
	Event       string          `json:"event"`
	IntentID    string          `json:"intent_id"`
	WorkspaceID string          `json:"workspace_id"`
	Amount      int64           `json:"amount"`
	Currency    string          `json:"currency"`
	Metadata    json.RawMessage `json:"metadata"`
}

const (
	EventPaymentSucceeded = "payment.succeeded"
	EventPaymentFailed    = "payment.failed"
	EventPaymentCancelled = "payment.cancelled"
)

type MarketplaceConfig struct {
	ID           string    `json:"id"`
	WorkspaceID  string    `json:"workspace_id"`
	CredentialID string    `json:"credential_id"`
	FeeBps       int       `json:"fee_bps"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SetMarketplaceConfigRequest struct {
	CredentialID string `json:"credential_id"`
	FeeBps       int    `json:"fee_bps"`
}
