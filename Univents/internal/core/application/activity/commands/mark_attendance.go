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

	ga := uc.gaClient

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

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("attendance").
		Action("mark").
		Scope(activity.ScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("activity").SetMessage("insufficient permissions")
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
