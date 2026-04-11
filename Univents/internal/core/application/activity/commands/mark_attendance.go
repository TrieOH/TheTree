package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

// FIXME Limit the times attendance can be marked to within the activity time

func (uc *CommandService) MarkAttendance(ctx context.Context, activityID, recordID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.MarkAttendance")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("mark.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var activity *domain.Activity
	activity, err = uc.activities.GetByID(ctx, activityID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("mark_attendance"),
		authz.Resource("activity", activity.ID.String()),
	); err != nil {
		return err
	}

	var attendanceRecord *domain.AttendanceRecord
	attendanceRecord, err = uc.activities.GetAttendanceRecordByID(ctx, recordID)
	if err != nil {
		return err
	}
	if attendanceRecord.Status != domain.AttendanceStatusRegistered {
		return errx.Invalid("attendance record").SetMessage("cannot mark attendance on activities on statuses different than registered")
	}

	if err = uc.activities.MarkAttendanceRecordStatus(ctx, recordID, &sub.ID, domain.AttendanceStatusCompleted); err != nil {
		return err
	}

	return nil
}
