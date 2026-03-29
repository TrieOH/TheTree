package types

import (
	"TrieForms/internal/shared/errx"
	"TrieForms/internal/shared/validation"
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID        uuid.UUID  `json:"id"`
	OwnerID   uuid.UUID  `json:"owner_id"`
	ScopeID   uuid.UUID  `json:"scope_id"`
	ProjectID uuid.UUID  `json:"project_id"`
	Name      string     `json:"name"`
	KeyHash   string     `json:"-"`      // bcrypt hash, never returned
	KeyPrefix string     `json:"prefix"` // first 8 chars for display e.g. "tf_live_"
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

func NewAPIKey(projectID, userID uuid.UUID, name, keyHash, keyPrefix string) (*APIKey, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("api_key").SetMessage("error generating uuid").SetCause(err)
	}

	ak := &APIKey{
		ID:        id,
		OwnerID:   userID,
		ProjectID: projectID,
		Name:      name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		CreatedAt: time.Now(),
	}

	if err := ak.validate(); err != nil {
		return nil, err
	}

	return ak, nil
}

func (k *APIKey) validate() error {
	return validation.Run(
		validation.RequireUUID("api_key", "owner_id", k.OwnerID),
		validation.RequireUUID("api_key", "project_id", k.ProjectID),
		validation.RequireString("api_key", "name", k.Name),
		validation.RequireString("api_key", "key_hash", k.KeyHash),
		validation.RequireString("api_key", "key_prefix", k.KeyPrefix),
	)
}

func (k *APIKey) AddScope(scopeID uuid.UUID) {
	k.ScopeID = scopeID
}
