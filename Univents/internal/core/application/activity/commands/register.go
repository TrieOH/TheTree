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

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("attend"),
		authz.Resource("activity", activity.ID.String()),
	); err != nil {
		return err
	}

	attendanceRecord := domain.NewAttendanceRecord(sub.ID, activity.ID)
	if _, err = uc.activities.Register(ctx, *attendanceRecord); err != nil {
		return err
	}

	return nil
}
