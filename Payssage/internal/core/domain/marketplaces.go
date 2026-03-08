package domain

import (
	"time"

	"github.com/google/uuid"
)

type MarketplaceConfig struct {
	ID           uuid.UUID
	WorkspaceID  uuid.UUID
	CredentialID uuid.UUID
	FeeBps       int // 500 = 5.00%
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
