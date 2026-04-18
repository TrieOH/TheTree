package authz

import (
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

type AuthMethod string

const (
	AuthMethodSession AuthMethod = "session"
	AuthMethodApiKey  AuthMethod = "api_key"
)

type Principal struct {
	UserID    uuid.UUID  `json:"user_id"`
	ProjectID *uuid.UUID `json:"project_id"`
	SessionID *uuid.UUID `json:"session_id"`
	Method    AuthMethod `json:"-"`
}

func NewPrincipal(
	ctx context.Context,
	access *contracts.AccessClaims,
) (*Principal, error) {
	if access == nil {
		return nil, fail.New(errx.TokenMissingAccessClaims).RecordCtx(ctx)
	}

	return &Principal{
		UserID:    access.Sub.ID,
		ProjectID: access.Sub.ProjectID,
		SessionID: &access.Sub.SessionID,
		Method:    AuthMethodSession,
	}, nil
}
