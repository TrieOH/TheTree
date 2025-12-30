package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionID uuid.UUID  `json:"session_id"`
	ProjectID *uuid.UUID `json:"project_id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenID   uuid.UUID  `json:"token_id"`
	IssuedAt  time.Time  `json:"issued_at"`
	UserAgent string     `json:"user_agent"`
	UserIp    string     `json:"user_ip"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	UserType  string     `json:"user_type"`
}

type SessionFilter struct {
	UserID        uuid.UUID
	SessionID     *uuid.UUID
	ExcludeID     *uuid.UUID
	TokenID       *uuid.UUID
	ExpiredBefore *time.Time
}
