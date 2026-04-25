package contracts

import (
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID        uuid.UUID  `json:"id"`
	OwnerID   uuid.UUID  `json:"owner_id"   validate:"required"`
	Name      string     `json:"name"       validate:"required"`
	KeyHash   string     `json:"-"          validate:"required"`
	KeyPrefix string     `json:"prefix"     validate:"required"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

func NewAPIKey(userID uuid.UUID, name, keyHash, keyPrefix string) (*APIKey, error) {
	ak := &APIKey{
		OwnerID:   userID,
		Name:      name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		CreatedAt: time.Now(),
	}
	if err := validate.Struct(ak); err != nil {
		return nil, err
	}
	return ak, nil
}
