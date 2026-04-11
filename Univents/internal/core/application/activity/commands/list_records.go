package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) ListRecords(ctx context.Context, activityID uuid.UUID) (records []domain.AttendanceRecord, err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.ListRecords")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("mark.success", err == nil))
	}()

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

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_attendance"),
		authz.Resource("activity", activity.ID.String()),
	); err != nil {
		return nil, err
	}

	var attendanceRecords []domain.AttendanceRecord
	attendanceRecords, err = uc.activities.ListActivityAttendanceRecords(ctx, activityID)
	if err != nil {
		return nil, err
	}

	return attendanceRecords, nil
}
