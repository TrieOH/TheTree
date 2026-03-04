package commands

import (
	"context"
	"errors"
	"time"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

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

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var event *domain.Event
	event, err = uc.events.GetEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("editions").
		Action("announce").
		Scope(event.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("edition").SetMessage("insufficient permissions")
	}

	edition, err := uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if edition.Status != domain.EditionStatusDraft {
		return errors.New("can't announce editions on statuses different than draft")
	}

	now := time.Now()
	if edition.StartsAt.Before(now.Add(time.Minute * 5)) {
		return errors.New("announcement must be at least 5 minutes out from right now")
	}

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
