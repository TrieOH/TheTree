package authz

import (
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/errx"
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
	UserID    uuid.UUID
	ProjectID *uuid.UUID
	Method    AuthMethod
}

func NewPrincipal(
	ctx context.Context,
	access *auth.AccessClaims,
	refresh *auth.RefreshClaims,
) (*Principal, error) {
	if access == nil {
		return nil, fail.New(errx.TokenMissingAccessClaims).RecordCtx(ctx)
	}
	if refresh == nil {
		return nil, fail.New(errx.TokenMissingRefreshClaims).RecordCtx(ctx)
	}

	return &Principal{
		UserID:    access.Sub.ID,
		ProjectID: access.Sub.ProjectID,
		Method:    AuthMethodSession,
	}, nil
}
