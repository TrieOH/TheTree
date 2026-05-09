package tickets

import (
	"context"
	"univents/internal/platform/database"
	"univents/internal/platform/database/sqlc"
	"univents/internal/shared/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"

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

var _ ports.TicketsRepository = (*ticketsRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.TicketsRepository {
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

func mapTicketFromDB(src *sqlc.Ticket) *contracts.Ticket {
	var dst contracts.Ticket
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

func mapTicketPermissionFromDB(src *sqlc.TicketPermission) *contracts.TicketPermission {
	return &contracts.TicketPermission{
		ID:             src.ID,
		TicketID:       src.TicketID,
		PermissionType: contracts.PermissionType(src.PermissionType),
		ActivityID:     src.ActivityID,
		ProductID:      src.ProductID,
		CheckpointID:   src.CheckpointID,
		CreatedAt:      src.CreatedAt,
	}
}

func (repo *ticketsRepo) Create(ctx context.Context, toCreate contracts.Ticket) (*contracts.Ticket, error) {
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

func (repo *ticketsRepo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Ticket, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketsRepo.GetByID")
	defer span.End()

	sqlcTicket, err := repo.queries(ctx).GetTicketByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "ticket")
	}

	return mapTicketFromDB(&sqlcTicket), nil
}

func (repo *ticketsRepo) AddPermission(ctx context.Context, toCreate contracts.TicketPermission) (*contracts.TicketPermission, error) {
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

func (repo *ticketsRepo) List(ctx context.Context, editionID uuid.UUID) ([]contracts.Ticket, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketsRepo.List")
	defer span.End()

	sqlcTickets, err := repo.queries(ctx).ListEditionTickets(ctx, editionID)
	if err != nil {
		return nil, errx.FromDB(err, "ticket")
	}

	out := make([]contracts.Ticket, 0, len(sqlcTickets))
	for _, ticket := range sqlcTickets {
		out = append(out, *mapTicketFromDB(&ticket))
	}
	return out, nil
}

func (repo *ticketsRepo) GetPermissions(ctx context.Context, ticketID uuid.UUID) ([]contracts.TicketPermission, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketsRepo.List")
	defer span.End()

	sqlcTicketPermissions, err := repo.queries(ctx).GetTicketPermissions(ctx, ticketID)
	if err != nil {
		return nil, errx.FromDB(err, "ticket")
	}

	out := make([]contracts.TicketPermission, 0, len(sqlcTicketPermissions))
	for _, perms := range sqlcTicketPermissions {
		out = append(out, *mapTicketPermissionFromDB(&perms))
	}
	return out, nil
}
