package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/plataform/telemetry"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (uc *CommandService) PublishEvent(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "EventService.PublishEvent")
	defer span.End()

	// FIXME send to BG Worker and return
	var auditor *domain.AuditBuilder
	defer func() {
		if auditor != nil {
			auditor.Emit()
			ae := uc.events.AppendEventAudits(ctx, auditor.GetAudits()) // FIXME make this outbox later
			if ae != nil {
				telemetry.Log().Error("failed to insert publish event audits", zap.Error(ae))
			}
		}
	}()

	ga := uc.gaClient

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	event, err := uc.events.GetEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	auditor = domain.StartAudit(event.ID, domain.ActorTypeUnknown, &sub.ID).
		Action(string(domain.EventAuditActionActivated)).
		State(domain.ActionStateFailed)

	isOwner := sub.ID == event.CreatedBy

	if isOwner {
		auditor.Actor(domain.ActorTypeOwner)
	}

	allowed, err := ga.Authz.Check().User(sub.ID).
		Object("events").
		Action("publish").
		Scope(event.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		auditor.AddMetadata("reason", "Internal System Error")
		return err
	}
	if !allowed {
		auditor.AddMetadata("reason", "Forbidden")
		return fail.New(errx.AuthzInsufficientPermissions)
	}
	if !isOwner {
		auditor.Actor(domain.ActorTypeAdmin)
	}

	if event.Status != domain.StatusDraft {
		auditor.AddMetadata("reason", "Invalid Status Change")
		return fail.New(errx.EventPublishNonDraft).WithArgs(string(event.Status))
	}

	err = uc.events.PublishEvent(ctx, eventID)
	if err != nil {
		auditor.AddMetadata("reason", err.Error())
		return err
	}

	auditor.State(domain.ActionStateSucceeded).StatusTo(string(domain.StatusActive))
	return nil
}
