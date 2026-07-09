package repos

import (
	"context"
	"lib/database"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) ListActivityAttendanceRecords(ctx context.Context, activityID uuid.UUID) ([]contracts.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.ListActivityAttendanceRecords")
	defer span.End()

	sqlcAttendanceRecords, err := database.Queries(ctx, repo.q).ListActivityAttendanceRecords(ctx, activityID)
	if err != nil {
		return nil, repo.dbe(err, "attendance record")
	}

	out := make([]contracts.AttendanceRecord, 0, len(sqlcAttendanceRecords))
	for _, record := range sqlcAttendanceRecords {
		out = append(out, *mapAttendanceRecordFromDB(&record))
	}
	return out, nil
}
