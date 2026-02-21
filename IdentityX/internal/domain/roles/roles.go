package roles

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID
	ProjectID   *uuid.UUID
	Name        string
	Description *string
	Meta        *json.RawMessage
	CreatedAt   time.Time
	UpdatedAt   time.Time

	ScopeID    *uuid.UUID
	ScopeName  *string
	ExternalID *string
}
