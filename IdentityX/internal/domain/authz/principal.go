package authz

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	"context"
	"encoding/json"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

type Principal struct {
	// ===== Identity =====
	UserID     uuid.UUID
	Email      string
	UserType   string
	ProjectID  *uuid.UUID
	Metadata   *json.RawMessage
	IsVerified bool
	VerifiedAt *time.Time

	// ===== Session =====
	SessionID uuid.UUID
	UserAgent string
	UserIP    string

	// ===== Token linkage =====
	AccessJTI  uuid.UUID
	RefreshJTI uuid.UUID

	// ===== Raw claims (escape hatch) =====
	AccessClaims  *auth.AccessClaims
	RefreshClaims *auth.RefreshClaims
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

	accessJTI, err := uuid.Parse(access.ID)
	if err != nil {
		return nil, fail.New(apierr.TokenAccessInvalidID).With(err).RecordCtx(ctx)
	}

	refreshJTI, err := uuid.Parse(refresh.ID)
	if err != nil {
		return nil, fail.New(apierr.TokenRefreshInvalidID).With(err).RecordCtx(ctx)
	}

	return &Principal{
		UserID:     access.Sub.ID,
		Email:      access.Sub.Email,
		UserType:   access.Sub.UserType,
		ProjectID:  access.Sub.ProjectID,
		Metadata:   access.Sub.Metadata,
		IsVerified: access.Sub.IsVerified,
		VerifiedAt: access.Sub.VerifiedAt,

		SessionID: access.Sub.SessionID,
		UserAgent: access.Sub.UserAgent,
		UserIP:    access.Sub.UserIP,

		AccessJTI:  accessJTI,
		RefreshJTI: refreshJTI,

		AccessClaims:  access,
		RefreshClaims: refresh,
	}, nil
}
