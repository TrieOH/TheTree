package project_users

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectUser struct {
	ID           uuid.UUID
	ProjectID    uuid.UUID
	Email        string
	PasswordHash string `json:"-"`
	UserType     string
	Metadata     *json.RawMessage
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time
}
