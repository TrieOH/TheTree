package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ApiKeys struct {
	ID         uuid.UUID       `json:"id"`
	ActorID    uuid.UUID       `json:"actor_id"`
	ProjectID  *uuid.UUID      `json:"project_id"`
	Name       string          `json:"name"`
	KeyPrefix  string          `json:"key_prefix"`
	KeyHash    string          `json:"key_hash"`
	Metadata   json.RawMessage `json:"metadata"`
	ExpiresAt  time.Time       `json:"expires_at"`
	RevokedAt  time.Time       `json:"revoked_at"`
	LastUsedAt time.Time       `json:"last_used_at"`
	CreatedAt  time.Time       `json:"created_at"`
}
