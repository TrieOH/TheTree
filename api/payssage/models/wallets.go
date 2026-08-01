package models

import (
	"time"

	"github.com/google/uuid"
)

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

func (w *Wallet) OwnedBy(id uuid.UUID) bool {
	return w.OwnerID == id
}

type CreateWalletInput struct {
	Name           string     `json:"name"`
	OrganizationID *uuid.UUID `json:"organization_id"`
}

type SetFeeBPSInput struct {
	WalletID       uuid.UUID  `json:"wallet_id"`
	FeeBps         int        `json:"fee_bps"`
	OrganizationID *uuid.UUID `json:"organization_id"`
}

type SetSandboxInput struct {
	WalletID       uuid.UUID  `json:"wallet_id"`
	Sandbox        bool       `json:"sandbox"`
	OrganizationID *uuid.UUID `json:"organization_id"`
}
