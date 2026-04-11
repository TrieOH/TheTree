package commands

import (
	"context"
	"univents/internal/commerce/domain"
	domain2 "univents/internal/core/domain"
	"univents/internal/shared/authz"

	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Create(ctx context.Context, in domain.CreateTicketSpec) (out *domain.Ticket, err error) {
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

	var validTicket *domain.Ticket
	validTicket, err = domain.NewTicket(sub.ID, in)
	if err != nil {
		return nil, err
	}

	var edition *domain2.Edition
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

	var created *domain.Ticket
	created, err = uc.tickets.Create(ctx, *validTicket)
	if err != nil {
		return nil, err
	}

	return created, nil
}
