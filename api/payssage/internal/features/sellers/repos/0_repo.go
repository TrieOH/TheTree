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

var _ ports.SellerRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.SellerRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("seller"),
	}
}

func mapSeller(src sqlc2.Seller) models.Seller {
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
