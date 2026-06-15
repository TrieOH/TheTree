package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BlacklistEntryType string

const (
	BlacklistEntryTypeActor  BlacklistEntryType = "actor"
	BlacklistEntryTypeToken  BlacklistEntryType = "token"
	BlacklistEntryTypeApiKey BlacklistEntryType = "api_key"
	BlacklistEntryTypeEmail  BlacklistEntryType = "email"
	BlacklistEntryTypeIP     BlacklistEntryType = "ip"
)

type BlacklistEntry struct {
	ID               uuid.UUID          `json:"id"`
	CreatedByActorID *uuid.UUID         `json:"created_by_actor_id"`
	ProjectID        *uuid.UUID         `json:"project_id"`
	Type             BlacklistEntryType `json:"type"`
	Target           string             `json:"target"`
	Reason           *string            `json:"reason"`
	Metadata         *json.RawMessage   `json:"metadata"`
	CreatedAt        time.Time          `json:"created_at"`
	ExpiresAt        *time.Time         `json:"expires_at"`
}
