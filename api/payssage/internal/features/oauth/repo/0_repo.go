package repo

import (
	"lib/database"
	"payssage/errs"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.OAuthStateRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrOAuthStateNotFound),
	}
}

func mapState(src sqlc.OauthState) models.OAuthState {
	return models.OAuthState{
		State:            src.State,
		WalletID:         src.WalletID,
		OrganizationID:   src.OrganizationID,
		OwnerID:          src.OwnerID,
		Provider:         src.Provider,
		Flow:             models.OAuthFlow(src.Flow),
		FinalRedirectURL: src.FinalRedirectUrl,
		CreatedAt:        src.CreatedAt,
		ExpiresAt:        src.ExpiresAt,
	}
}
