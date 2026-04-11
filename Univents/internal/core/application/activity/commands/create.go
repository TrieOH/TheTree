package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"

	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Create(ctx context.Context, in domain.CreateActivitySpec) (out *domain.Activity, err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var edition *domain.Edition
	edition, err = uc.editions.GetByID(ctx, in.EditionID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_activity"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return nil, err
	}

	var validActivity *domain.Activity
	validActivity, err = domain.NewActivity(sub.ID, in, edition)
	if err != nil {
		return nil, err
	}

	var created *domain.Activity
	created, err = uc.activities.Create(ctx, validActivity)
	if err != nil {
		return nil, err
	}

	return created, nil
}
