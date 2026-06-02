package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"Informd/ports"
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

var _ ports.ResponseRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ResponseRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("response"),
	}
}

func mapResponse(src sqlc.Response) models.Response {
	return models.Response{
		ID:          src.ID,
		FormID:      src.FormID,
		InviteID:    src.InviteID,
		ResponderID: src.ResponderID,
		Email:       src.Email,
		StartedAt:   src.StartedAt,
		FinishedAt:  src.FinishedAt,
	}
}
