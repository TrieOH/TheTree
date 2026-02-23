package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/plataform/telemetry"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (uc *CommandService) CreateEvent(ctx context.Context, in domain.CreateEventSpec) (out *domain.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	// FIXME send to BG Worker and return
	var auditor *domain.AuditBuilder
	defer func() {
		if auditor != nil {
			auditor.Emit()
			ae := uc.events.AppendEventAudits(ctx, auditor.GetAudits()) // FIXME make this outbox later
			if ae != nil {
				telemetry.Log().Error("failed to insert create event audits", zap.Error(ae))
			}
		}
	}()

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var validEvent *domain.Event
	validEvent, err = domain.NewEvent(sub.ID, &sub.ID, in)
	if err != nil {
		return nil, fail.AsFail(err).System().RecordCtx(ctx)
	}

	auditor = domain.StartAudit(validEvent.ID, domain.ActorTypeUnknown, &sub.ID).
		Action(string(domain.EventAuditActionCreated)).
		State(domain.ActionStateFailed)

	isOwner := sub.ID == validEvent.CreatedBy

	if isOwner {
		auditor.Actor(domain.ActorTypeOwner)
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("events").
		Action("create").
		Allowed(ctx)
	if err != nil {
		auditor.AddMetadata("reason", "Internal System Error")
		return nil, fail.AsFail(err).System().RecordCtx(ctx)
	}
	if !allowed {
		auditor.AddMetadata("reason", "Forbidden")
		return nil, fail.New(errx.AuthzInsufficientPermissions)
	}

	if !isOwner {
		auditor.Actor(domain.ActorTypeAdmin)
	}

	span.SetAttributes(attribute.String("event.id", validEvent.ID.String()))

	var scope *goauth.Scope
	var idStr = validEvent.ID.String()
	scope, err = ga.Scopes.Create(ctx, validEvent.Slug, &idStr)
	if err != nil {
		auditor.AddMetadata("reason", "Internal System Error")
		return nil, fail.AsFail(err).System().RecordCtx(ctx)
	}
	validEvent.AddScope(scope.ID)

	roleID, err := uuid.Parse("019c8138-bc3b-7bfa-9140-7d5fe5e3bc73")
	if err != nil {
		auditor.AddMetadata("reason", "Internal System Error")
		telemetry.Log().Error("Error parsing UUID", zap.Error(err))
		return nil, err
	}

	err = uc.gaClient.Roles.GiveToUser(ctx, sub.ID, roleID, &scope.ID)
	if err != nil {
		auditor.AddMetadata("reason", "Internal System Error")
		return nil, fail.AsFail(err).System().RecordCtx(ctx)
	}

	var created *domain.Event
	created, err = uc.events.CreateEvent(ctx, validEvent) // FIXME if this fails the scope must be undone (SAGA PATTERN)
	if err != nil {
		auditor.AddMetadata("reason", err.Error())
		return nil, err
	}

	auditor.State(domain.ActionStateSucceeded).StatusTo(string(domain.StatusDraft))
	return created, nil
}
