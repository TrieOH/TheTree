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

var _ ports.CollectorRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.CollectorRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("collector"),
	}
}

func mapCollector(src sqlc.Collector) models.Collector {
	return models.Collector{
		ID:             src.ID,
		OwnerID:        src.OwnerID,
		OrganizationID: src.OrganizationID,
		Provider:       src.Provider,
		ProviderUserID: src.ProviderUserID,
		Credentials:    src.Credentials,
		CreatedAt:      src.CreatedAt,
		RevokedAt:      src.RevokedAt,
	}
}
