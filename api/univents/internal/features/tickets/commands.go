package tickets

import (
	"context"

	"lib/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	editions ports.EditionsRepository
	tickets  ports.TicketsRepository
	asynq    *asynq.Client
	tracer   trace.Tracer
	az       *authzed.Client
	tx       database.TxRunner
}

func NewCommandService(
	editions ports.EditionsRepository,
	tickets ports.TicketsRepository,
	asynq *asynq.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		editions: editions,
		tickets:  tickets,
		asynq:    asynq,
		tracer:   tracer,
		tx:       tx,
	}
}

func (uc *CommandService) Create(ctx context.Context, in contracts.CreateTicketSpec) (out *contracts.Ticket, err error) {
	ctx, span := uc.tracer.Start(ctx, "TicketService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var validTicket *contracts.Ticket
	validTicket, err = contracts.NewTicket(sub.ID, in)
	if err != nil {
		return nil, err
	}

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, in.EditionID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_tickets"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return nil, err
	}

	var created *contracts.Ticket
	created, err = uc.tickets.Create(ctx, *validTicket)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *CommandService) AddPermission(ctx context.Context, in contracts.CreateTicketPermissionSpec) (out *contracts.TicketPermission, err error) {
	ctx, span := uc.tracer.Start(ctx, "TicketService.AddPermission")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("add.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var validTicketPermission *contracts.TicketPermission
	validTicketPermission, err = contracts.NewTicketPermission(sub.ID, in)
	if err != nil {
		return nil, err
	}

	var ticket *contracts.Ticket
	ticket, err = uc.tickets.GetByID(ctx, in.TicketID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("ticket", ticket.ID.String()),
	); err != nil {
		return nil, err
	}

	var created *contracts.TicketPermission
	created, err = uc.tickets.AddPermission(ctx, *validTicketPermission)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *CommandService) RemovePermission(ctx context.Context, id, ticketID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "TicketService.RemovePermission")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("add.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var ticket *contracts.Ticket
	ticket, err = uc.tickets.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("ticket", ticket.ID.String()),
	); err != nil {
		return err
	}

	err = uc.tickets.RemovePermission(ctx, id, ticketID)
	if err != nil {
		return err
	}

	return nil
}
