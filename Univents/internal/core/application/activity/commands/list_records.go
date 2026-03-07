package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) ListRecords(ctx context.Context, activityID uuid.UUID) (records []domain.AttendanceRecord, err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.ListRecords")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("mark.success", err == nil))
	}()

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var activity *domain.Activity
	activity, err = uc.activities.GetByID(ctx, activityID)
	if err != nil {
		return nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("activities").
		Action("manage").
		Scope(activity.ScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("activity").SetMessage("insufficient permissions")
	}

	var attendanceRecords []domain.AttendanceRecord
	attendanceRecords, err = uc.activities.ListActivityAttendanceRecords(ctx, activityID)
	if err != nil {
		return nil, err
	}

	return attendanceRecords, nil
}
