package commands

import (
	"context"
	"fmt"
	"io"
	"time"
	"univents/internal/core/domain"
	"univents/internal/plataform/telemetry"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GoAuthUserClient is the assumed interface for bulk user fetches from GoAuth.
// ga.Users.GetBatch() is expected to satisfy this.
type GoAuthUserClient interface {
	GetBatch(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]domain.ExportUserInfo, error)
}

// exportRow is the flat struct the csvwriter reflects over.
// Field order doesn't matter — the frontend's selected key order controls the CSV column order.
// Add a new column: add a field with a csv tag, done.
type exportRow struct {
	// User fields — populated from GoAuth batch fetch
	UserID   uuid.UUID `csv:"user_id"`
	FullName string    `csv:"full_name"`
	Email    string    `csv:"email"`

	// Activity fields
	ActivityID       uuid.UUID               `csv:"activity_id"`
	ActivityTitle    string                  `csv:"activity_title"`
	ActivityLocation string                  `csv:"activity_location"`
	ActivityStartsAt time.Time               `csv:"activity_starts_at"`
	ActivityEndsAt   time.Time               `csv:"activity_ends_at"`
	ActivityStatus   domain.ActivityStatus   `csv:"activity_status"`
	Difficulty       *domain.DifficultyLevel `csv:"difficulty"`

	// Attendance fields
	AttendanceID     uuid.UUID               `csv:"attendance_id"`
	AttendanceStatus domain.AttendanceStatus `csv:"attendance_status"`
	CheckedInAt      *time.Time              `csv:"checked_in_at"`
	CancelledAt      *time.Time              `csv:"cancelled_at"`
	RegisteredAt     time.Time               `csv:"registered_at"`

	// Scan metadata
	ScannedBy *uuid.UUID `csv:"scanned_by"`
}

// userKeys is the set of column keys that require a GoAuth batch fetch.
// If none of these appear in the request we skip the network call entirely.
var userKeys = map[string]struct{}{
	"user_id":   {},
	"full_name": {},
	"email":     {},
}

func (uc *CommandService) ExportAttendance(ctx context.Context, editionID uuid.UUID, req domain.ExportRequest, w io.Writer) error {
	ctx, span := uc.tracer.Start(ctx, "ExportAttendanceCSV.ExportAttendance")
	defer span.End()

	// TODO: permission check — verify caller has export rights for this edition

	dbRows, err := uc.activities.AttendanceExport(ctx, editionID, req.Filters)
	if err != nil {
		telemetry.Log().Error("attendance export query failed",
			zap.String("edition_id", editionID.String()),
			zap.Any("filters", req.Filters),
			zap.Error(err),
		)
		return fmt.Errorf("querying export rows: %w", err)
	}

	telemetry.Log().Info("attendance export query succeeded", zap.Int("rows", len(dbRows)))

	var userMap map[uuid.UUID]domain.ExportUserInfo
	if needsUserData(req.Columns) && len(dbRows) > 0 {
		//userIDs := collectUserIDs(dbRows)
		//userMap, err = uc.goauth.GetBatch(ctx, userIDs)
		//if err != nil {
		//	return fmt.Errorf("fetching users from goauth: %w", err)
		//}
	}

	rows := toExportRows(dbRows, userMap)

	return uc.csvWriter.Stream(w, req.Columns, rows)
}

func needsUserData(cols []string) bool {
	for _, col := range cols {
		if _, ok := userKeys[col]; ok {
			return true
		}
	}
	return false
}

func toExportRows(dbRows []domain.AttendanceExportRow, userMap map[uuid.UUID]domain.ExportUserInfo) []exportRow {
	rows := make([]exportRow, len(dbRows))
	for i, r := range dbRows {
		user := userMap[r.UserID]
		rows[i] = exportRow{
			UserID:           r.UserID,
			FullName:         user.FullName,
			Email:            user.Email,
			ActivityID:       r.ActivityID,
			ActivityTitle:    r.ActivityTitle,
			ActivityLocation: r.ActivityLocation,
			ActivityStartsAt: r.ActivityStartsAt,
			ActivityEndsAt:   r.ActivityEndsAt,
			ActivityStatus:   r.ActivityStatus,
			Difficulty:       r.Difficulty,
			AttendanceID:     r.AttendanceID,
			AttendanceStatus: r.AttendanceStatus,
			CheckedInAt:      r.CheckedInAt,
			CancelledAt:      r.CancelledAt,
			RegisteredAt:     r.RegisteredAt,
			ScannedBy:        r.ScannedBy,
		}
	}
	return rows
}
