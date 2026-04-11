package queries

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"

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

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_editions"),
		authz.Resource("event", eventID.String()),
	); err != nil {
		return nil, err
	}

	var outEditions []domain.Edition
	outEditions, err = uc.editions.ListAdmin(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return outEditions, nil
}
