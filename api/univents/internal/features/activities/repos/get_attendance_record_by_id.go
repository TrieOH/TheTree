package repos

import (
	"context"
	"lib/database"
	"univents/contracts"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) GetAttendanceRecordByID(ctx context.Context, id uuid.UUID) (*contracts.AttendanceRecord, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.GetAttendanceRecordByID")
	defer span.End()

	sqlcAttendanceRecord, err := database.Queries(ctx, repo.q).GetAttendanceRecordByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err, "attendance record")
	}

	return mapAttendanceRecordFromDB(&sqlcAttendanceRecord), nil
}
