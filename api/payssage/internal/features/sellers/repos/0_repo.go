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

var _ ports.SellerRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries, tracer trace.Tracer) *Repo {
	return &Repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("seller"),
	}
}

func mapSeller(src sqlc.Seller) models.Seller {
	return models.Seller{
		ID:             src.ID,
		WalletID:       src.WalletID,
		Provider:       src.Provider,
		ProviderUserID: src.ProviderUserID,
		Credentials:    src.Credentials,
		CreatedAt:      src.CreatedAt,
		RevokedAt:      src.RevokedAt,
	}
}
