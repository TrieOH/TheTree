package queries

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

func (uc *QueryService) ListEventAudits(ctx context.Context, eventID uuid.UUID) (out []domain.Audit, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EventService.ListEventAudits")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.gaClient.Authz.Check().User(sub.ID).
		Object("events").
		Action("administrate").
		Scope(eventID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fail.New(errx.AuthzInsufficientPermissions).RecordCtx(ctx)
	}

	var outAudits []domain.Audit
	outAudits, err = uc.events.ListEventAuditByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return outAudits, nil
}
