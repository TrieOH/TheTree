package project

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID       `json:"id"`
	ProjectName string          `json:"project_name"`
	OwnerID     uuid.UUID       `json:"owner_id"`
	Metadata    json.RawMessage `json:"metadata"`
	IsActive    bool            `json:"is_active"`
	PubKey      string          `json:"-"`
	PrivKey     []byte          `json:"-"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
