package models

import (
	"time"

	"github.com/google/uuid"
)

type RevokedRefreshToken struct {
	TokenID   uuid.UUID `json:"token_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
