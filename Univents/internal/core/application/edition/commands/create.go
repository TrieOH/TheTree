package commands

import (
	"context"
	"encoding/json"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Create(ctx context.Context, in domain.CreateEditionSpec) (out *domain.Edition, err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.Create")
	defer span.End()

	if err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		out, err = uc.createInternal(ctx, in)
		return err
	}); err != nil {
		return &domain.Edition{}, err
	}

	return out, nil
}

func (uc *CommandService) createInternal(ctx context.Context, in domain.CreateEditionSpec) (out *domain.Edition, err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.createInternal")
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

	var validEdition *domain.Edition
	validEdition, err = domain.NewEdition(sub.ID, in)
	if err != nil {
		return nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("editions").
		Action("create").
		Scope(in.GoAuthEventScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("edition").SetMessage("insufficient permissions")
	}

	span.SetAttributes(attribute.String("event.id", validEdition.ID.String()))

	meta := json.RawMessage(`{"color": "#a84bfa", "icon": "Tickets"}`)
	var scope *goauth.Scope
	var idStr = validEdition.ID.String()
	scope, err = ga.Scopes.CreateWithParent(ctx, validEdition.EditionName, &idStr, &in.GoAuthEventScopeID, meta)
	if err != nil {
		return nil, err
	}
	validEdition.AddScope(scope.ID)

	var created *domain.Edition
	created, err = uc.editions.Create(ctx, validEdition) // FIXME if this fails the scope must be undone (SAGA PATTERN)
	if err != nil {
		return nil, err
	}

	err = uc.events.AddEdition(ctx, validEdition.EventID)
	if err != nil {
		return nil, err
	}

	return created, nil
}
