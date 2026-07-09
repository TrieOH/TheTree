package commands

import (
	"context"
	"errors"
	"time"
	"univents/contracts"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) Announce(ctx context.Context, editionID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.Announce")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("announce.success", err == nil))
	}()

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if edition.Status != contracts.EditionStatusDraft {
		return errors.New("can't announce editions on statuses different than draft")
	}

	//var task *asynq.Task
	opensAt := edition.RegistrationOpensAt
	if opensAt == nil {
		opensAt = new(time.Now())
	}
	//task, err = contracts.NewOpenEditionTask(edition.ID, *opensAt)
	//if err != nil {
	//	return err
	//}
	//if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
	//	return err
	//}
	_ = uc.editions.Open(ctx, edition.ID)

	//task, err = contracts.NewStartEditionTask(edition.ID, edition.StartsAt)
	//if err != nil {
	//	return err
	//}
	//if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
	//	return err
	//}
	_ = uc.editions.Start(ctx, edition.ID)

	//task, err = contracts.NewFinishEditionTask(edition.ID, edition.EndsAt)
	//if err != nil {
	//	return err
	//}
	//if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
	//	return err
	//}

	if err = uc.editions.Announce(ctx, editionID); err != nil {
		return err
	}

	return nil
}
