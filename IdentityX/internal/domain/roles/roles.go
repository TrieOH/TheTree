package roles

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID
	ProjectID   *uuid.UUID
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
