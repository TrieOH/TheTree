package commands

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Create(ctx context.Context, in domain.CreateTicketSpec) (out *domain.Ticket, err error) {
	ctx, span := uc.tracer.Start(ctx, "TicketService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var validTicket *domain.Ticket
	validTicket, err = domain.NewTicket(sub.ID, in)
	if err != nil {
		return nil, fail.AsFail(err).System().RecordCtx(ctx)
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("tickets").
		Action("create").
		Scope(in.EditionScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, fail.AsFail(err).System().RecordCtx(ctx)
	}
	if !allowed {
		return nil, fail.New(errx.AuthzInsufficientPermissions)
	}

	span.SetAttributes(attribute.String("ticket.id", validTicket.ID.String()))

	var created *domain.Ticket
	created, err = uc.tickets.Create(ctx, *validTicket)
	if err != nil {
		return nil, err
	}

	return created, nil
}
