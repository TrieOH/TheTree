package apierr

import (
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/authz"

	"github.com/MintzyG/fail"
	"github.com/google/uuid"
)

func NewPrincipal(
	access *auth.AccessClaims,
	refresh *auth.RefreshClaims,
) (*authz.Principal, error) {
	if access == nil {
		return nil, fail.New(TokenMissingAccessClaims)
	}
	if refresh == nil {
		return nil, fail.New(TokenMissingRefreshClaims)
	}

	accessJTI, err := uuid.Parse(access.ID)
	if err != nil {
		return nil, fail.New(TokenAccessInvalidID).With(err)
	}

	refreshJTI, err := uuid.Parse(refresh.ID)
	if err != nil {
		return nil, fail.New(TokenRefreshInvalidID).With(err)
	}

	return &authz.Principal{
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
