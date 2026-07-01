package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.CapabilityRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.CapabilityRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("capabilities"),
	}
}

func mapCapability(src sqlc.Capability) models.Capability {
	return models.Capability{
		ID:        src.ID,
		ProjectID: src.ProjectID,
		Resource:  src.Resource,
		Action:    src.Action,
		CreatedBy: src.CreatedBy,
		CreatedAt: src.CreatedAt,
	}
}
