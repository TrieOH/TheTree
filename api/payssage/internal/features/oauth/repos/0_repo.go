package repos

import (
	"lib/database"
	sqlc2 "payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
)

type repo struct {
	q      *sqlc2.Queries
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.OAuthStateRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.OAuthStateRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("state"),
	}
}

func mapState(src sqlc2.OauthState) models.OAuthState {
	return models.OAuthState{
		State:               src.State,
		WalletID:            src.WalletID,
		OrganizationID:      src.OrganizationID,
		OwnerID:             src.OwnerID,
		Provider:            src.Provider,
		Flow:                models.OAuthFlow(src.Flow),
		FinalRedirectUrl:    src.FinalRedirectUrl,
		ProviderRedirectUrl: src.ProviderRedirectUrl,
		CreatedAt:           src.CreatedAt,
		ExpiresAt:           src.ExpiresAt,
	}
}
