package infrastructure

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
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
	dst.ScopeID = src.ScopeID
	dst.EditionID = src.EditionID
	dst.Name = src.Name
	dst.Description = src.Description
	dst.CreatedBy = src.CreatedBy
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.DeletedAt = src.DeletedAt
	return &dst
}

func mapTicketPermissionFromDB(src *sqlc.TicketPermission) *domain.TicketPermission {
	return &domain.TicketPermission{
		ID:             src.ID,
		TicketID:       src.TicketID,
		PermissionType: domain.PermissionType(src.PermissionType),
		ActivityID:     src.ActivityID,
		ProductID:      src.ProductID,
		CheckpointID:   src.CheckpointID,
		CreatedAt:      src.CreatedAt,
	}
}

func (repo *ticketsRepo) Create(ctx context.Context, toCreate domain.Ticket) (*domain.Ticket, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketsRepo.Create")
	defer span.End()

	sqlcTicket, err := repo.queries(ctx).CreateTicket(ctx, sqlc.CreateTicketParams{
		ID:          toCreate.ID,
		EditionID:   toCreate.EditionID,
		ScopeID:     toCreate.ScopeID,
		Name:        toCreate.Name,
		Description: toCreate.Description,
		CreatedBy:   toCreate.CreatedBy,
	})
	if err != nil {
		return nil, errx.FromDB(err, "ticket")
	}

	return mapTicketFromDB(&sqlcTicket), nil
}

func (repo *ticketsRepo) AddPermission(ctx context.Context, toCreate domain.TicketPermission) (*domain.TicketPermission, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketsRepo.AddPermission")
	defer span.End()

	sqlcTicketPermission, err := repo.queries(ctx).AddTicketPermission(ctx, sqlc.AddTicketPermissionParams{
		ID:             toCreate.ID,
		TicketID:       toCreate.TicketID,
		PermissionType: sqlc.PermissionType(toCreate.PermissionType),
		ActivityID:     toCreate.ActivityID,
		ProductID:      toCreate.ProductID,
		CheckpointID:   toCreate.CheckpointID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "ticket permission")
	}

	return mapTicketPermissionFromDB(&sqlcTicketPermission), nil
}

func (repo *ticketsRepo) RemovePermission(ctx context.Context, id, ticketID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "TicketsRepo.RemovePermission")
	defer span.End()

	err := repo.queries(ctx).RemoveTicketPermission(ctx, sqlc.RemoveTicketPermissionParams{
		ID:       id,
		TicketID: ticketID,
	})
	if err != nil {
		return errx.FromDB(err, "ticket permission")
	}

	return nil
}
