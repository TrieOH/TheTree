package commands

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

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

	var ticket *domain.Ticket
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
