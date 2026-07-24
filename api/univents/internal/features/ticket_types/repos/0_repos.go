package repos

import (
	"univents/internal/database/sqlc"
	"univents/models"
	"univents/ports"

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

var _ ports.TicketTypeRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.TicketTypeRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("ticket_type"),
	}
}

func mapTicketType(src sqlc.TicketType) models.TicketType {
	return models.TicketType{
		ID:          src.ID,
		EditionID:   src.EditionID,
		Name:        src.Name,
		Description: src.Description,
		AccessLevel: src.AccessLevel,
		PriceCents:  src.Price,
		MaxQuantity: src.MaxQuantity,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}
