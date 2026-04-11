package commands

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"

	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) AddPermission(ctx context.Context, in domain.CreateTicketPermissionSpec) (out *domain.TicketPermission, err error) {
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

	var validTicketPermission *domain.TicketPermission
	validTicketPermission, err = domain.NewTicketPermission(sub.ID, in)
	if err != nil {
		return nil, err
	}

	var ticket *domain.Ticket
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

	var created *domain.TicketPermission
	created, err = uc.tickets.AddPermission(ctx, *validTicketPermission)
	if err != nil {
		return nil, err
	}

	return created, nil
}
