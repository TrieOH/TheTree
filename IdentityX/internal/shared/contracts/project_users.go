package contracts

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectUser struct {
	ID           uuid.UUID        `json:"id"`
	ProjectID    uuid.UUID        `json:"project_id"`
	Email        string           `json:"email"`
	PasswordHash string           `json:"-"`
	UserType     string           `json:"user_type"`
	Metadata     *json.RawMessage `json:"metadata"`
	IsActive     bool             `json:"is_active"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	LastLoginAt  *time.Time       `json:"last_login_at"`
	IsVerified   bool             `json:"is_verified"`
	VerifiedAt   *time.Time       `json:"verified_at"`
}
