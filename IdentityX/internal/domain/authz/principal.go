package authz

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

type Principal struct {
	UserID     uuid.UUID
	UserType   string
	ProjectID  *uuid.UUID
	IsVerified bool
}

func NewPrincipal(
	ctx context.Context,
	access *auth.AccessClaims,
	refresh *auth.RefreshClaims,
) (*Principal, error) {
	if access == nil {
		return nil, fail.New(apierr.TokenMissingAccessClaims).RecordCtx(ctx)
	}
	if refresh == nil {
		return nil, fail.New(apierr.TokenMissingRefreshClaims).RecordCtx(ctx)
	}

	return &Principal{
		UserID:     access.Sub.ID,
		UserType:   access.Sub.UserType,
		ProjectID:  access.Sub.ProjectID,
		IsVerified: access.Sub.IsVerified,
	}, nil
}
