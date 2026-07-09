package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/contracts"

	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) Create(ctx context.Context, in contracts.CreateActivitySpec) (out *contracts.Activity, err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, in.EditionID)
	if err != nil {
		return nil, err
	}

	var validActivity *contracts.Activity
	validActivity, err = contracts.NewActivity(ident.Sub.ID, in, edition)
	if err != nil {
		return nil, err
	}

	var created *contracts.Activity
	created, err = uc.activities.Create(ctx, validActivity)
	if err != nil {
		return nil, err
	}

	return created, nil
}
