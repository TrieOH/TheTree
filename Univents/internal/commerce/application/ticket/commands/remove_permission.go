package commands

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) RemovePermission(ctx context.Context, id, ticketID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "TicketService.RemovePermission")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("add.success", err == nil))
	}()

	ga := uc.gaClient

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

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("tickets").
		Action("edit").
		Scope(ticket.ScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("ticket permission").SetMessage("insufficient permissions")
	}

	span.SetAttributes(attribute.String("ticket_permission.id", id.String()))

	err = uc.tickets.RemovePermission(ctx, id, ticketID)
	if err != nil {
		return err
	}

	return nil
}
