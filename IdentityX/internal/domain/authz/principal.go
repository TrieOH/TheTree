package authz

import (
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/errx"
	"context"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

type AuthMethod string

const (
	AuthMethodSession        AuthMethod = "session"
	AuthMethodApiKey         AuthMethod = "api_key"
	AuthMethodServiceSession AuthMethod = "service_session"
)

type Principal struct {
	UserID    uuid.UUID
	ProjectID *uuid.UUID
	Method    AuthMethod
}

type ServiceSnapshot struct {
	UserID    uuid.UUID
	ProjectID *uuid.UUID
	UserType  string
	GlobalSID uuid.UUID
	AccessJTI string
	Issuer    string
	ExpiresAt time.Time
}

func (ss ServiceSnapshot) ToPrincipal() *Principal {
	return &Principal{
		UserID:    ss.UserID,
		ProjectID: ss.ProjectID,
		Method:    AuthMethodServiceSession,
	}
}

func NewPrincipal(
	ctx context.Context,
	access *auth.AccessClaims,
) (*Principal, error) {
	if access == nil {
		return nil, fail.New(errx.TokenMissingAccessClaims).RecordCtx(ctx)
	}

	return &Principal{
		UserID:    access.Sub.ID,
		ProjectID: access.Sub.ProjectID,
		Method:    AuthMethodSession,
	}, nil
}
