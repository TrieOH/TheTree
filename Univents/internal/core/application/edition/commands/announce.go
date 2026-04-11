package commands

import (
	"context"
	"errors"
	"time"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Announce(ctx context.Context, eventID, editionID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.Announce")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("announce.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var edition *domain.Edition
	edition, err = uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("announce"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return err
	}

	if edition.Status != domain.EditionStatusDraft {
		return errors.New("can't announce editions on statuses different than draft")
	}

	now := time.Now()
	var task *asynq.Task
	opensAt := edition.RegistrationOpensAt
	if opensAt == nil {
		opensAt = &now
	}
	task, err = domain.NewOpenEditionTask(edition.ID, *opensAt)
	if err != nil {
		return err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		return err
	}

	task, err = domain.NewStartEditionTask(edition.ID, edition.StartsAt)
	if err != nil {
		return err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		return err
	}

	task, err = domain.NewFinishEditionTask(edition.ID, edition.EndsAt)
	if err != nil {
		return err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		return err
	}

	if err = uc.editions.Announce(ctx, editionID); err != nil {
		return err
	}

	return nil
}
