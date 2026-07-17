package repos

import (
	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/models"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.OAuthStateRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.OAuthStateRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("state"),
	}
}

func mapState(src sqlc.OauthState) models.OAuthState {
	return models.OAuthState{
		State:               src.State,
		WalletID:            src.WalletID,
		Provider:            src.Provider,
		Flow:                models.OAuthFlow(src.Flow),
		FinalRedirectUrl:    src.FinalRedirectUrl,
		ProviderRedirectUrl: src.ProviderRedirectUrl,
		CreatedAt:           src.CreatedAt,
		ExpiresAt:           src.ExpiresAt,
	}
}
