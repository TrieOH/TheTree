package repos

import (
	"context"
	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/models"
)

func (repo *repo) Create(ctx context.Context, state models.OAuthState) (*models.OAuthState, error) {
	ctx, span := repo.tracer.Start(ctx, "Create")
	defer span.End()
	sqlcState, err := database.Queries(ctx, repo.q).CreateOAuthState(ctx, sqlc.CreateOAuthStateParams{
		State:            state.State,
		WalletID:         state.WalletID,
		Provider:         state.Provider,
		Flow:             state.Flow.String(),
		FinalRedirectUrl: state.FinalRedirectUrl,
		ExpiresAt:        state.ExpiresAt,
	})
	return new(mapState(sqlcState)), repo.dbe(err)
}
