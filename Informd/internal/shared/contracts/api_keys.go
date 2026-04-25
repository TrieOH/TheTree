package contracts

import (
	"time"

	"github.com/MintzyG/FastUtilitiesNet"
	"github.com/google/uuid"
)

type APIKey struct {
	ID        uuid.UUID  `json:"id"`
	OwnerID   uuid.UUID  `json:"owner_id"   validate:"required"`
	ProjectID uuid.UUID  `json:"project_id" validate:"required"`
	Name      string     `json:"name"       validate:"required"`
	KeyHash   string     `json:"-"          validate:"required"`
	KeyPrefix string     `json:"prefix"     validate:"required"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

func NewAPIKey(projectID, userID uuid.UUID, name, keyHash, keyPrefix string) (*APIKey, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fun.Errf("error generating uuid for api key: %s", err.Error()).Internal()
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

	if err = validate.Struct(ak); err != nil {
		return nil, err
	}

	return ak, nil
}
