package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/contracts"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) Unregister(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Unregister")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("unregister.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var isRegistered bool
	isRegistered, err = uc.activities.IsRegistered(ctx, ident.Sub.ID, activity.ID)
	if err != nil {
		return err
	}
	if !isRegistered {
		return errx.Invalid("activity").SetMessage("user isn't registered")
	}

	if err = uc.activities.Unregister(ctx, ident.Sub.ID, activity.ID); err != nil {
		return err
	}

	return nil
}
