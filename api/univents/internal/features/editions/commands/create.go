package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/contracts"

	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) Create(ctx context.Context, in contracts.CreateEditionSpec) (out *contracts.Edition, err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.Create")
	defer span.End()

	if err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		out, err = uc.createInternal(ctx, in)
		return err
	}); err != nil {
		return &contracts.Edition{}, err
	}

	return out, nil
}

func (uc *Commands) createInternal(ctx context.Context, in contracts.CreateEditionSpec) (out *contracts.Edition, err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.createInternal")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var validEdition *contracts.Edition
	validEdition, err = contracts.NewEdition(ident.Sub.ID, in)
	if err != nil {
		return nil, err
	}

	_, err = uc.events.GetByID(ctx, in.EventID)
	if err != nil {
		return nil, err
	}

	var created *contracts.Edition
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
