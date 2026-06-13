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

var _ ports.BlacklistRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.BlacklistRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("blacklist entry"),
	}
}

func mapEntry(src sqlc.BlacklistEntry) models.BlacklistEntry {
	return models.BlacklistEntry{
		ID:               src.ID,
		CreatedByActorID: src.CreatedByActorID,
		ProjectID:        src.ProjectID,
		Type:             models.BlacklistEntryType(src.Type),
		Target:           src.Target,
		Reason:           src.Reason,
		Metadata:         src.Metadata,
		CreatedAt:        src.CreatedAt,
		ExpiresAt:        src.ExpiresAt,
	}
}
