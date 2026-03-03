package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Create(ctx context.Context, in domain.CreateActivitySpec) (out *domain.Activity, err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	ga := uc.gaClient

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

	var validActivity *domain.Activity
	validActivity, err = domain.NewActivity(sub.ID, in, edition)
	if err != nil {
		return nil, fail.AsFail(err).System().RecordCtx(ctx)
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("activities").
		Action("create").
		Scope(in.EditionScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, fail.AsFail(err).System().RecordCtx(ctx)
	}
	if !allowed {
		return nil, fail.New(errx.AuthzInsufficientPermissions)
	}

	span.SetAttributes(attribute.String("activity.id", validActivity.ID.String()))

	var scope *goauth.Scope
	var idStr = validActivity.ID.String()
	scope, err = ga.Scopes.CreateWithParent(ctx, validActivity.Title, &idStr, &in.EditionScopeID)
	if err != nil {
		return nil, fail.AsFail(err).System().RecordCtx(ctx)
	}
	validActivity.AddScope(scope.ID)

	var created *domain.Activity
	created, err = uc.activities.Create(ctx, validActivity)
	if err != nil {
		return nil, err
	}

	return created, nil
}
