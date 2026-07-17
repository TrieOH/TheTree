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

type CreateWalletRequest struct {
	Name           string     `json:"name"`
	OrganizationID *uuid.UUID `json:"organization_id"`
}

func (r CreateWalletRequest) ToInput() CreateWalletInput {
	return CreateWalletInput{
		Name:           r.Name,
		OrganizationID: r.OrganizationID,
	}
}

type CreateWalletInput struct {
	Name           string     `json:"name"`
	OrganizationID *uuid.UUID `json:"organization_id"`
}

type SetFeeBPSRequest struct {
	FeeBps         int        `json:"fee_bps"`
	OrganizationID *uuid.UUID `json:"organization_id"`
}

func (r SetFeeBPSRequest) ToInput(walletID uuid.UUID) SetFeeBPSInput {
	return SetFeeBPSInput{
		WalletID:       walletID,
		FeeBps:         r.FeeBps,
		OrganizationID: r.OrganizationID,
	}
}

type SetFeeBPSInput struct {
	WalletID       uuid.UUID  `json:"wallet_id"`
	FeeBps         int        `json:"fee_bps"`
	OrganizationID *uuid.UUID `json:"organization_id"`
}

type SetSandboxRequest struct {
	Sandbox        bool       `json:"sandbox"`
	OrganizationID *uuid.UUID `json:"organization_id"`
}

func (r SetSandboxRequest) ToInput(walletID uuid.UUID) SetSandboxInput {
	return SetSandboxInput{
		WalletID:       walletID,
		Sandbox:        r.Sandbox,
		OrganizationID: r.OrganizationID,
	}
}

type SetSandboxInput struct {
	WalletID       uuid.UUID  `json:"wallet_id"`
	Sandbox        bool       `json:"sandbox"`
	OrganizationID *uuid.UUID `json:"organization_id"`
}
