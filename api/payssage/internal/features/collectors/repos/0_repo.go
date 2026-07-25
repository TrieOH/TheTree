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

var _ ports.CollectorRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.CollectorRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("collector"),
	}
}

func mapCollector(src sqlc2.Collector) models.Collector {
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
