package inbounds

import (
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/session"
	"encoding/json"
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

type PrincipalOutput struct {
	UserID        uuid.UUID
	Email         string
	UserType      string
	ProjectID     *uuid.UUID
	Metadata      *json.RawMessage
	SessionID     uuid.UUID
	UserAgent     string
	UserIP        string
	AccessJTI     uuid.UUID
	RefreshJTI    uuid.UUID
	AccessClaims  *auth.AccessClaims
	RefreshClaims *auth.RefreshClaims
}

func PrincipalToPrincipalOutput(p authz.Principal) *PrincipalOutput {
	return &PrincipalOutput{
		UserID:        p.UserID,
		Email:         p.Email,
		UserType:      p.UserType,
		ProjectID:     p.ProjectID,
		Metadata:      p.Metadata,
		SessionID:     p.SessionID,
		UserAgent:     p.UserAgent,
		UserIP:        p.UserIP,
		AccessJTI:     p.AccessJTI,
		RefreshJTI:    p.RefreshJTI,
		AccessClaims:  p.AccessClaims,
		RefreshClaims: p.RefreshClaims,
	}
}

type ErrRevokeCurrentSession struct{}

func (e ErrRevokeCurrentSession) Error() string {
	return "cannot revoke the currently active session"
}

type ErrSessionNotFound struct{}

func (e ErrSessionNotFound) Error() string {
	return "session not found or revoked"
}

type ErrSessionUnauthorized struct{}

func (e ErrSessionUnauthorized) Error() string {
	return "session not found or revoked"
}
