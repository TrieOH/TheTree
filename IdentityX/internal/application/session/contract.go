package session

import (
	"GoAuth/internal/domain/session"
	"time"

	"github.com/google/uuid"
)

type OutputSession struct {
	SessionID uuid.UUID
	ProjectID *uuid.UUID
	UserID    uuid.UUID
	IssuedAt  time.Time
	UserAgent string
	UserIp    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	UserType  string
}

func OutputSessionFromSession(s *session.Session) *OutputSession {
	return &OutputSession{
		SessionID: s.SessionID,
		ProjectID: s.ProjectID,
		UserID:    s.UserID,
		IssuedAt:  s.IssuedAt,
		UserAgent: s.UserAgent,
		UserIp:    s.UserIP,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		UserType:  s.UserType,
	}
}

func OutputSessionSliceFromSessionSlice(src []session.Session) []OutputSession {
	dst := make([]OutputSession, 0, len(src))
	for _, s := range src {
		dst = append(dst, *OutputSessionFromSession(&s))
	}
	return dst
}
