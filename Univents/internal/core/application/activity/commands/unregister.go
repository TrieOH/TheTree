package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Unregister(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Unregister")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("unregister.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var activity *domain.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var isRegistered bool
	isRegistered, err = uc.activities.IsRegistered(ctx, sub.ID, activity.ID)
	if err != nil {
		return err
	}
	if !isRegistered {
		return errx.Invalid("activity").SetMessage("user isn't registered")
	}

	if err = uc.activities.Unregister(ctx, sub.ID, activity.ID); err != nil {
		return err
	}

	return nil
}
