package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"

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

	var event *domain.Event
	event, err = uc.events.GetByID(ctx, in.EventID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_editions"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
		return nil, err
	}

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
