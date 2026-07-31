package repos

import (
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
)

type Repo struct {
	q      *sqlc.Queries
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.OAuthStateRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries, tracer trace.Tracer) *Repo {
	return &Repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("state"),
	}
}

func mapState(src sqlc.OauthState) models.OAuthState {
	return models.OAuthState{
		State:               src.State,
		WalletID:            src.WalletID,
		OrganizationID:      src.OrganizationID,
		OwnerID:             src.OwnerID,
		Provider:            src.Provider,
		Flow:                models.OAuthFlow(src.Flow),
		FinalRedirectURL:    src.FinalRedirectUrl,
		ProviderRedirectURL: src.ProviderRedirectUrl,
		CreatedAt:           src.CreatedAt,
		ExpiresAt:           src.ExpiresAt,
	}
}
