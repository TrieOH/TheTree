package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *activitiesRepo) MarkAttendanceRecordStatus(ctx context.Context, id uuid.UUID, scannedBy *uuid.UUID, status contracts.AttendanceStatus) error {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.MarkAttendanceRecordStatus")
	defer span.End()

	err := database.Queries(ctx, repo.q).MarkAttendanceRecordStatus(ctx, sqlc.MarkAttendanceRecordStatusParams{
		Status:    sqlc.AttendanceStatus(status),
		ScannedBy: scannedBy,
		ID:        id,
	})
	if err != nil {
		return repo.dbe(err, "attendance record")
	}

	return nil
}
