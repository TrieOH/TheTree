package queries

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
)

func (uc *QueryService) List(ctx context.Context, editionID uuid.UUID) (out []domain.Activity, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ActivityService.List")
	defer span.End()

	return uc.activities.List(ctx, editionID)
}

func (uc *QueryService) AdminList(ctx context.Context, editionID uuid.UUID) (out []domain.Activity, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ActivityService.AdminList")
	defer span.End()

	edition, err := uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	allowed, err := ga.Authz.Check().
		User(sub.ID).
		Object("activities").
		Action("read").
		Scope(edition.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("activities").SetMessage("insufficient permissions")
	}

	return uc.activities.ListAdmin(ctx, editionID)
}
