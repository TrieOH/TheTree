package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionID uuid.UUID
	ProjectID *uuid.UUID
	UserID    uuid.UUID
	TokenID   uuid.UUID
	IssuedAt  time.Time
	UserAgent string
	UserIp    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	UserType  string
}

type Filter struct {
	UserID        uuid.UUID
	SessionID     *uuid.UUID
	ExcludeID     *uuid.UUID
	TokenID       *uuid.UUID
	ExpiredBefore *time.Time
}
