package dto

import (
	"GoAuth/internal/application/session"
	"time"

	"github.com/google/uuid"
)

type SessionResponse struct {
	SessionID uuid.UUID  `json:"session_id"`
	ProjectID *uuid.UUID `json:"project_id"`
	UserID    uuid.UUID  `json:"user_id"`
	IssuedAt  time.Time  `json:"issued_at"`
	UserAgent string     `json:"user_agent"`
	UserIp    string     `json:"user_ip"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	UserType  string     `json:"user_type"`
}

func SessionResponseFromSessionOutput(s session.OutputSession) SessionResponse {
	return SessionResponse{
		SessionID: s.SessionID,
		ProjectID: s.ProjectID,
		UserID:    s.UserID,
		IssuedAt:  s.IssuedAt,
		UserAgent: s.UserAgent,
		UserIp:    s.UserIp,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		UserType:  s.UserType,
	}
}

func SessionResponseSliceFromSessionOutputSlice(src []session.OutputSession) []SessionResponse {
	dst := make([]SessionResponse, 0, len(src))
	for _, s := range src {
		dst = append(dst, SessionResponseFromSessionOutput(s))
	}
	return dst
}
