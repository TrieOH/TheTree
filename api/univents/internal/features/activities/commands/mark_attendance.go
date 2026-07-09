package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/contracts"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

// FIXME Limit the times attendance can be marked to within the activity time

func (uc *Commands) MarkAttendance(ctx context.Context, activityID, recordID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.MarkAttendance")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("mark.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	_, err = uc.activities.GetByID(ctx, activityID)
	if err != nil {
		return err
	}

	var attendanceRecord *contracts.AttendanceRecord
	attendanceRecord, err = uc.activities.GetAttendanceRecordByID(ctx, recordID)
	if err != nil {
		return err
	}
	if attendanceRecord.Status != contracts.AttendanceStatusRegistered {
		return errx.Invalid("attendance record").SetMessage("cannot mark attendance on activities on statuses different than registered")
	}

	if err = uc.activities.MarkAttendanceRecordStatus(ctx, recordID, &ident.Sub.ID, contracts.AttendanceStatusCompleted); err != nil {
		return err
	}

	return nil
}
