package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	ID         uuid.UUID       `json:"id"`
	ActorID    uuid.UUID       `json:"actor_id"`
	ProjectID  *uuid.UUID      `json:"project_id"`
	Name       string          `json:"name"`
	KeyPrefix  string          `json:"key_prefix"`
	KeyHash    string          `json:"key_hash"`
	Metadata   json.RawMessage `json:"metadata"`
	ExpiresAt  *time.Time      `json:"expires_at"`
	RevokedAt  *time.Time      `json:"revoked_at"`
	LastUsedAt *time.Time      `json:"last_used_at"`
	CreatedAt  time.Time       `json:"created_at"`
}

type CreateApiKeyRequest struct {
	Name                    string     `json:"name"`
	ExpiresAt               *time.Time `json:"expires_at"`
	CreateForServiceAccount bool       `json:"create_for_service_account"`
}

func (r CreateApiKeyRequest) ToInput(projectID uuid.UUID) CreateApiKeyInput {
	return CreateApiKeyInput{
		Name:                    r.Name,
		ExpiresAt:               r.ExpiresAt,
		ProjectID:               projectID,
		CreateForServiceAccount: r.CreateForServiceAccount,
	}
}

type CreateApiKeyInput struct {
	Name                    string
	ExpiresAt               *time.Time
	ProjectID               uuid.UUID
	CreateForServiceAccount bool
}

type CreateApiKeyResponse struct {
	Key    *ApiKey `json:"key"`
	RawKey string  `json:"raw_key"`
}
