package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	ID            uuid.UUID       `json:"id"`
	SubjectID     uuid.UUID       `json:"subject_id"`
	Name          string          `json:"name"`
	DisplayPrefix string          `json:"display_prefix"`
	KeyHash       []byte          `json:"key_hash"`
	Metadata      json.RawMessage `json:"metadata"`
	ExpiresAt     *time.Time      `json:"expires_at"`
	RevokedAt     *time.Time      `json:"revoked_at"`
	LastUsedAt    *time.Time      `json:"last_used_at"`
	CreatedBy     uuid.UUID       `json:"created_by"`
	CreatedAt     time.Time       `json:"created_at"`
}

type CreateApiKeyRequest struct {
	SubjectID    *uuid.UUID  `json:"subject_id"`
	Capabilities []uuid.UUID `json:"capabilities"`
	Name         string      `json:"name"`
	Env          string      `json:"env"`
	ExpiresAt    *time.Time  `json:"expires_at"`
}

func (r CreateApiKeyRequest) ToInput(projectID *uuid.UUID) CreateApiKeyInput {
	return CreateApiKeyInput{
		SubjectID:    r.SubjectID,
		Capabilities: r.Capabilities,
		Name:         r.Name,
		Env:          r.Env,
		ExpiresAt:    r.ExpiresAt,
		ProjectID:    projectID,
	}
}

type CreateApiKeyInput struct {
	SubjectID    *uuid.UUID  `json:"subject_id"`
	Capabilities []uuid.UUID `json:"capabilities"`
	Name         string      `json:"name"`
	Env          string      `json:"env"`
	ExpiresAt    *time.Time  `json:"expires_at"`
	ProjectID    *uuid.UUID
}

type CreateApiKeyResponse struct {
	Key    *ApiKey `json:"key"`
	RawKey string  `json:"raw_key"`
}
