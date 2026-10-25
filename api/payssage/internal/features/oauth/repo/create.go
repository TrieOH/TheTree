package repo

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *Repo) Create(ctx context.Context, state models.OAuthState) (*models.OAuthState, error) {
	sqlcState, err := database.Queries(ctx, repo.q).CreateOAuthState(ctx, sqlc.CreateOAuthStateParams{
		State:            state.State,
		WalletID:         state.WalletID,
		OrganizationID:   state.OrganizationID,
		OwnerID:          state.OwnerID,
		Provider:         state.Provider,
		Flow:             state.Flow.String(),
		FinalRedirectUrl: state.FinalRedirectURL,
		ExpiresAt:        state.ExpiresAt,
	})
	return new(mapState(sqlcState)), repo.dbe(err)
}
