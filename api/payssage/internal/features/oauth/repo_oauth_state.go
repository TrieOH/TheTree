package oauth

import (
	"context"
	"payssage/ports"

	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/internal/shared/errx"
	"payssage/models"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type oauthStatesRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.OAuthStateRepo = (*oauthStatesRepo)(nil)

func NewOAuthStatesRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.OAuthStateRepo {
	return &oauthStatesRepo{q: q, log: log, tracer: tracer}
}

func (repo *oauthStatesRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapOAuthStateFromDB(src *sqlc.OauthState) *models.OAuthState {
	return &models.OAuthState{
		State:            src.State,
		WorkspaceID:      src.WorkspaceID,
		Provider:         src.Provider,
		Flow:             src.Flow,
		IsMarketplace:    src.IsMarketplace,
		FeeBps:           src.FeeBps,
		FinalRedirectURL: src.FinalRedirectUrl,
		CreatedAt:        src.CreatedAt,
		ExpiresAt:        src.ExpiresAt,
	}
}

func (repo *oauthStatesRepo) Create(ctx context.Context, state models.OAuthState) (*models.OAuthState, error) {
	ctx, span := repo.tracer.Start(ctx, "OAuthStatesRepo.Create")
	defer span.End()

	row, err := repo.queries(ctx).CreateOAuthState(ctx, sqlc.CreateOAuthStateParams{
		State:            state.State,
		WorkspaceID:      state.WorkspaceID,
		Provider:         state.Provider,
		Flow:             state.Flow,
		IsMarketplace:    state.IsMarketplace,
		FeeBps:           state.FeeBps,
		FinalRedirectUrl: state.FinalRedirectURL,
		ExpiresAt:        state.ExpiresAt,
	})
	if err != nil {
		return nil, errx.FromDB(err, "oauth_state")
	}

	return mapOAuthStateFromDB(&row), nil
}

func (repo *oauthStatesRepo) Get(ctx context.Context, state string) (*models.OAuthState, error) {
	ctx, span := repo.tracer.Start(ctx, "OAuthStatesRepo.Get")
	defer span.End()

	row, err := repo.queries(ctx).GetOAuthState(ctx, state)
	if err != nil {
		return nil, errx.FromDB(err, "oauth_state")
	}

	return mapOAuthStateFromDB(&row), nil
}

func (repo *oauthStatesRepo) Delete(ctx context.Context, state string) error {
	ctx, span := repo.tracer.Start(ctx, "OAuthStatesRepo.Delete")
	defer span.End()

	if err := repo.queries(ctx).DeleteOAuthState(ctx, state); err != nil {
		return errx.FromDB(err, "oauth_state")
	}

	return nil
}
