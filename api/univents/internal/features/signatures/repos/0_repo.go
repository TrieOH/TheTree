package repos

import (
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"
	"univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.SignatureRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.SignatureRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("signatures"),
	}
}

func mapSignature(src sqlc.Signature) models.Signature {
	return models.Signature{
		ID:        src.ID,
		EditionID: src.EditionID,
		Title:     src.Title,
		URL:       src.Url,
	}
}
