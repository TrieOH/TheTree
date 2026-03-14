package project

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID
	ProjectName string
	Domain      string
	OwnerID     uuid.UUID
	Metadata    json.RawMessage
	IsActive    bool
	PubKey      string
	PrivKey     []byte `json:"-"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
