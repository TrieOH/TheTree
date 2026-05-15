package contracts

import (
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	ProjectID uuid.UUID `json:"project_id"`
	ClientID  uuid.UUID `json:"client_id"`
	KeyHash   string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
