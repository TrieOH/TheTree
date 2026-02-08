package inbounds

import (
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/session"
	"time"

	"github.com/google/uuid"
)

type OutputSession struct {
	SessionID  uuid.UUID
	ProjectID  *uuid.UUID
	FamilyID   uuid.UUID
	IdentityID uuid.UUID
	IssuedAt   time.Time
	UserAgent  string
	UserIp     string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserType   string
}

func OutputSessionSliceFromSessionSlice(src []session.Session) []OutputSession {
	dst := make([]OutputSession, 0, len(src))
	for i := range src {
		dst = append(dst, *OutputSessionFromSession(&src[i]))
	}
	return dst
}

func OutputSessionFromSession(s *session.Session) *OutputSession {
	return &OutputSession{
		SessionID:  s.SessionID,
		FamilyID:   s.FamilyID,
		ProjectID:  s.ProjectID,
		IdentityID: s.IdentityID,
		IssuedAt:   s.IssuedAt,
		UserAgent:  s.UserAgent,
		UserIp:     s.UserIP,
		ExpiresAt:  s.ExpiresAt,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
		UserType:   s.UserType,
	}
}

type MeOutput struct {
	AccessClaims      *auth.AccessClaims
	RefreshExpireDate time.Time
}
