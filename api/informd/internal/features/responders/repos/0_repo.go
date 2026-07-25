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

var _ ports.ResponderRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.ResponderRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("responder"),
	}
}

func mapResponder(src sqlc2.Responder) models.Responder {
	return models.Responder{
		ID:     src.ID,
		UserID: src.UserID,
		Email:  src.Email,
	}
}
