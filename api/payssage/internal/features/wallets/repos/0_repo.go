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

var _ ports.WalletRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.WalletRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("wallet"),
	}
}

func mapWallet(src sqlc.Wallet) models.Wallet {
	return models.Wallet{
		ID:             src.ID,
		OwnerID:        src.OwnerID,
		OrganizationID: src.OrganizationID,
		Name:           src.Name,
		Sandbox:        src.Sandbox,
		FeeBps:         src.FeeBps,
		CollectorID:    src.CollectorID,
		CreatedAt:      src.CreatedAt,
	}
}
