package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Register(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Register")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("register.success", err == nil))
	}()

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var activity *domain.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var isRegistered bool
	isRegistered, err = uc.activities.IsRegistered(ctx, sub.ID, activity.ID)
	if err != nil {
		return err
	}
	if isRegistered {
		return errx.Invalid("activity").SetMessage("user already registered to activity")
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("activities").
		Action("attend").
		Scope(activity.ScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("activity").SetMessage("insufficient permissions")
	}

	attendanceRecord := domain.NewAttendanceRecord(sub.ID, activity.ID)
	if _, err = uc.activities.Register(ctx, *attendanceRecord); err != nil {
		return err
	}

	return nil
}
