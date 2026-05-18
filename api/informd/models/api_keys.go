package models

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
	}
	if err := validate.Struct(ak); err != nil {
		return nil, err
	}
	return ak, nil
}

type CreateAPIKeyRequest struct {
	Name string `json:"name" validate:"required"`
}

type APIKeyResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Prefix    string     `json:"prefix"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

type CreateAPIKeyResponse struct {
	APIKeyResponse
	Key string `json:"key"` // only returned once
}

func (ak *APIKey) ToResponse() APIKeyResponse {
	return APIKeyResponse{
		ID:        ak.ID,
		Name:      ak.Name,
		Prefix:    ak.KeyPrefix,
		CreatedAt: ak.CreatedAt,
		RevokedAt: ak.RevokedAt,
	}
}
