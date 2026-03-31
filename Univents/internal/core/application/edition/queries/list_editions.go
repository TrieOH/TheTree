package queries

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
)

func (uc *QueryService) ListEditions(ctx context.Context, eventID uuid.UUID) (out []domain.Edition, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EditionsService.ListEditions")
	defer span.End()

	var outEditions []domain.Edition
	outEditions, err = uc.editions.List(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return outEditions, nil
}

func (uc *QueryService) ListEditionsAdmin(ctx context.Context, eventID uuid.UUID) (out []domain.Edition, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EditionsService.ListEditionsAdmin")
	defer span.End()

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var event *domain.Event
	event, err = uc.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("editions").
		Action("read").
		Scope(event.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("edition").SetMessage("insufficient permissions")
	}

	var outEditions []domain.Edition
	outEditions, err = uc.editions.ListAdmin(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return outEditions, nil
}
