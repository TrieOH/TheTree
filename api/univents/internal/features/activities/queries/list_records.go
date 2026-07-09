package queries

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *Queries) ListRecords(ctx context.Context, activityID uuid.UUID) (records []contracts.AttendanceRecord, err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.ListRecords")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("mark.success", err == nil))
	}()

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, activityID)
	if err != nil {
		return nil, err
	}

	var attendanceRecords []contracts.AttendanceRecord
	attendanceRecords, err = uc.activities.ListActivityAttendanceRecords(ctx, activity.ID)
	if err != nil {
		return nil, err
	}

	return attendanceRecords, nil
}
