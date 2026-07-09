package commands

import (
	"context"
	"errors"
	"univents/contracts"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) Publish(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Publish")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("publish.success", err == nil))
	}()

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if activity.Status != contracts.ActivityStatusDraft {
		return errors.New("can't publish activities on statuses different than draft")
	}

	//var task *asynq.Task
	//task, err = contracts.NewStartActivityTask(activity.ID, activity.StartsAt)
	//if err != nil {
	//	return err
	//}
	//if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
	//	return err
	//}
	_ = uc.activities.Start(ctx, activity.ID)

	//task, err = contracts.NewEndActivityTask(activity.ID, activity.EndsAt)
	//if err != nil {
	//	return err
	//}
	//if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
	//	return err
	//}

	if err = uc.activities.Publish(ctx, activity.ID); err != nil {
		return err
	}

	return nil
}
