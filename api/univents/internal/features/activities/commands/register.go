package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/contracts"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) Register(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Register")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("register.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var isRegistered bool
	isRegistered, err = uc.activities.IsRegistered(ctx, ident.Sub.ID, activity.ID)
	if err != nil {
		return err
	}
	if isRegistered {
		return errx.Invalid("activity").SetMessage("user already registered to activity")
	}

	attendanceRecord := contracts.NewAttendanceRecord(ident.Sub.ID, activity.ID)
	if _, err = uc.activities.Register(ctx, *attendanceRecord); err != nil {
		return err
	}

	return nil
}
