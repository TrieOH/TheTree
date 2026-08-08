package payssage

import (
	"time"

	"github.com/google/uuid"
)

// Wallet mirrors `models.Wallet` — the Payssage ownership unit for sellers,
// webhook endpoints, and payment intents.
type Wallet struct {
	ID             uuid.UUID  `json:"id"`
	OwnerID        uuid.UUID  `json:"owner_id"`
	OrganizationID *uuid.UUID `json:"organization_id"`
	Name           string     `json:"name"`
	Sandbox        bool       `json:"sandbox"`
	FeeBps         int        `json:"fee_bps"`
	CollectorID    *uuid.UUID `json:"collector_id"`
	CreatedAt      time.Time  `json:"created_at"`
}

// Seller mirrors `models.Seller` — a provider account (e.g. MercadoPago)
// bound to a wallet, created by the provider OAuth flow.
type Seller struct {
	ID             uuid.UUID  `json:"id"`
	WalletID       uuid.UUID  `json:"wallet_id"`
	Provider       string     `json:"provider"`
	ProviderUserID string     `json:"provider_user_id"`
	CreatedAt      time.Time  `json:"created_at"`
	RevokedAt      *time.Time `json:"revoked_at"`
}
