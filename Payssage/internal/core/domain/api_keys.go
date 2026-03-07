package domain

import (
	"TriePayments/internal/shared/errx"
	"TriePayments/internal/shared/validation"
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID          uuid.UUID  `json:"id"`
	ScopeID     uuid.UUID  `json:"scope_id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	Name        string     `json:"name"`
	KeyHash     string     `json:"-"`      // bcrypt hash, never returned
	KeyPrefix   string     `json:"prefix"` // first 8 chars for display e.g. "tp_live_"
	CreatedAt   time.Time  `json:"created_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
}

func NewAPIKey(workspaceID uuid.UUID, name, keyHash, keyPrefix string) (*APIKey, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("product").SetMessage("error generating uuid").SetCause(err)
	}

	ak := &APIKey{
		ID:          id,
		WorkspaceID: workspaceID,
		Name:        name,
		KeyHash:     keyHash,
		KeyPrefix:   keyPrefix,
		CreatedAt:   time.Now(),
	}

	if err := ak.validate(); err != nil {
		return nil, err
	}

	return ak, nil
}

func (k *APIKey) validate() error {
	return validation.Run(
		validation.RequireUUID("api_key", "workspace_id", k.WorkspaceID),
		validation.RequireString("api_key", "name", k.Name),
		validation.RequireString("api_key", "key_hash", k.KeyHash),
		validation.RequireString("api_key", "key_prefix", k.KeyPrefix),
	)
}

func (k *APIKey) AddScope(scopeID uuid.UUID) {
	k.ScopeID = scopeID
}
