package commands

import (
	"context"
	"errors"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Publish(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Publish")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("publish.success", err == nil))
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

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("publish"),
		authz.Resource("activity", activity.ID.String()),
	); err != nil {
		return err
	}

	if activity.Status != domain.ActivityStatusDraft {
		return errors.New("can't publish activities on statuses different than draft")
	}

	var task *asynq.Task
	task, err = domain.NewStartActivityTask(activity.ID, activity.StartsAt)
	if err != nil {
		return err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		return err
	}

	task, err = domain.NewEndActivityTask(activity.ID, activity.EndsAt)
	if err != nil {
		return err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		return err
	}

	if err = uc.activities.Publish(ctx, activity.ID); err != nil {
		return err
	}

	return nil
}
