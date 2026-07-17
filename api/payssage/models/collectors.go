package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Collector struct {
	ID             uuid.UUID       `json:"id"`
	OwnerID        uuid.UUID       `json:"owner_id"`
	OrganizationID *uuid.UUID      `json:"organization_id"`
	Provider       string          `json:"provider"`
	ProviderUserID string          `json:"provider_user_id"`
	Credentials    json.RawMessage `json:"credentials"`
	CreatedAt      time.Time       `json:"created_at"`
	RevokedAt      *time.Time      `json:"revoked_at"`
}
