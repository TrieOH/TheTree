package contracts

import (
	"time"

	"github.com/google/uuid"
)

type MarketplaceConfig struct {
	ID           uuid.UUID `json:"id"`
	WorkspaceID  uuid.UUID `json:"workspace_id"`
	CredentialID uuid.UUID `json:"credential_id"`
	Provider     string    `json:"provider"`
	FeeBps       int       `json:"fee_bps"` // 500 = 5.00%
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
