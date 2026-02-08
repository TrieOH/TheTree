package apikey

import (
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	ProjectID uuid.UUID
	ClientID  uuid.UUID
	KeyHash   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
