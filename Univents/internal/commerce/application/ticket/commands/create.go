package commands

import (
	"context"
	"encoding/json"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/TrieOH/goauth-sdk-go"
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
		return nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("tickets").
		Action("create").
		Scope(in.EditionScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("ticket").SetMessage("insufficient permissions")
	}

	span.SetAttributes(attribute.String("ticket.id", validTicket.ID.String()))

	meta := json.RawMessage(`{"color": "#aa21ff", "icon": "TicketCheck", "folder": "tickets"}`)
	var scope *goauth.Scope
	var idStr = validTicket.ID.String()
	scope, err = ga.Scopes.CreateWithParent(ctx, validTicket.Name, &idStr, &in.EditionScopeID, meta)
	if err != nil {
		return nil, err
	}
	validTicket.AddScope(scope.ID)

	var created *domain.Ticket
	created, err = uc.tickets.Create(ctx, *validTicket)
	if err != nil {
		return nil, err
	}

	return created, nil
}
