package authz

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	"encoding/json"

	"github.com/google/uuid"
)

type Principal struct {
	// ===== Identity =====
	UserID    uuid.UUID
	Email     string
	UserType  string
	ProjectID *uuid.UUID
	Metadata  *json.RawMessage

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
	access *auth.AccessClaims,
	refresh *auth.RefreshClaims,
) (*Principal, error) {

	accessJTI, err := uuid.Parse(access.ID)
	if err != nil {
		return nil, apierr.ErrUnauthorized.WithMsg("couldn't parse access JTI").WithID(apierr.TokenInvalidID).WithCause(err)
	}

	refreshJTI, err := uuid.Parse(refresh.ID)
	if err != nil {
		return nil, apierr.ErrUnauthorized.WithMsg("couldn't parse refresh JTI").WithID(apierr.TokenInvalidID).WithCause(err)

	}

	return &Principal{
		UserID:    access.Sub.ID,
		Email:     access.Sub.Email,
		UserType:  access.Sub.UserType,
		ProjectID: access.Sub.ProjectID,
		Metadata:  access.Sub.Metadata,

		SessionID: access.Sub.SessionID,
		UserAgent: access.Sub.UserAgent,
		UserIP:    access.Sub.UserIP,

		AccessJTI:  accessJTI,
		RefreshJTI: refreshJTI,

		AccessClaims:  access,
		RefreshClaims: refresh,
	}, nil
}
