package queries

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"

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

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_activities"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return nil, err
	}

	return uc.activities.ListAdmin(ctx, editionID)
}
