package repos

import (
	sqlc2 "Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type repo struct {
	q      *sqlc2.Queries
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.ResponseRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.ResponseRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("response"),
	}
}

func mapResponse(src sqlc2.Response) models.Response {
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
