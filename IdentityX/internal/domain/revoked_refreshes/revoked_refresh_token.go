package revoked_refreshes

import (
	"time"

	"github.com/google/uuid"
)

type RevokedRefreshToken struct {
	TokenID   uuid.UUID
	CreatedAt time.Time
	ExpiresAt time.Time
}
