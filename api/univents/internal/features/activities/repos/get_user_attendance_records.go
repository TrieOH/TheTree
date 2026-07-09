package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) GetUserActivityAttendanceRecords(ctx context.Context, userID, activityID uuid.UUID) ([]contracts.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetUserActivityAttendanceRecords")
	defer span.End()

	sqlcRecords, err := database.Queries(ctx, repo.q).GetUserActivityAttendanceRecords(ctx, sqlc.GetUserActivityAttendanceRecordsParams{
		ActivityID: activityID,
		UserID:     userID,
	})
	if err != nil {
		return nil, repo.dbe(err, "attendance records")
	}

	out := make([]contracts.AttendanceRecord, 0, len(sqlcRecords))
	for _, record := range sqlcRecords {
		out = append(out, *mapAttendanceRecordFromDB(&record))
	}
	return out, nil
}
