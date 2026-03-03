package infrastructure

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"

	"github.com/MintzyG/fail/v3"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ticketsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ domain.TicketsRepository = (*ticketsRepo)(nil)

func NewTicketsRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) domain.TicketsRepository {
	return &ticketsRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *ticketsRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapTicketFromDB(src *sqlc.Ticket) *domain.Ticket {
	var dst domain.Ticket
	dst.ID = src.ID
	dst.EditionID = src.EditionID
	dst.Name = src.Name
	dst.Description = src.Description
	dst.CreatedBy = src.CreatedBy
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.DeletedAt = src.DeletedAt
	return &dst
}

func (repo *ticketsRepo) Create(ctx context.Context, toCreate domain.Ticket) (*domain.Ticket, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketsRepo.Create")
	defer span.End()

	sqlcTicket, err := repo.queries(ctx).CreateTicket(ctx, sqlc.CreateTicketParams{
		EditionID:   toCreate.EditionID,
		Name:        toCreate.Name,
		Description: toCreate.Description,
		CreatedBy:   toCreate.CreatedBy,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return mapTicketFromDB(&sqlcTicket), nil
}
