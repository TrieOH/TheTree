package domain

import (
	"time"

	"github.com/google/uuid"
)

// ExportFilters defines the optional filters applied before generating the export.
// All fields are optional — zero value means no filter applied.
type ExportFilters struct {
	ActivityIDs        []uuid.UUID        `json:"activity_ids"`
	AttendanceStatuses []AttendanceStatus `json:"attendance_statuses"`
	ActivityStatuses   []ActivityStatus   `json:"activity_statuses"`
	Difficulties       []DifficultyLevel  `json:"difficulties"`
	DateFrom           *time.Time         `json:"date_from"`
	DateTo             *time.Time         `json:"date_to"`
}

// ExportRequest is the full payload sent by the frontend.
// Columns is an ordered list of csv tag keys matching the exportRow struct fields.
// The csvwriter validates keys and controls column order.
type ExportRequest struct {
	Columns []string      `json:"columns"`
	Filters ExportFilters `json:"filters"`
}

// AttendanceExportRow is the flat struct the csvwriter reflects over.
// Add a new export column: add a field with a csv tag, done.
type AttendanceExportRow struct {
	// User fields — populated from GoAuth batch fetch
	UserID   uuid.UUID `csv:"user_id"`
	FullName string    `csv:"full_name"`
	Email    string    `csv:"email"`

	// Activity fields
	ActivityID       uuid.UUID        `csv:"activity_id"`
	ActivityTitle    string           `csv:"activity_title"`
	ActivityLocation string           `csv:"activity_location"`
	ActivityStartsAt time.Time        `csv:"activity_starts_at"`
	ActivityEndsAt   time.Time        `csv:"activity_ends_at"`
	ActivityStatus   ActivityStatus   `csv:"activity_status"`
	Difficulty       *DifficultyLevel `csv:"difficulty"`

	// Attendance fields
	AttendanceID     uuid.UUID        `csv:"attendance_id"`
	AttendanceStatus AttendanceStatus `csv:"attendance_status"`
	CheckedInAt      *time.Time       `csv:"checked_in_at"`
	CancelledAt      *time.Time       `csv:"cancelled_at"`
	RegisteredAt     time.Time        `csv:"registered_at"`

	// Scan metadata
	ScannedBy *uuid.UUID `csv:"scanned_by"`
}

// ExportUserInfo holds the GoAuth-sourced user data merged into each row.
type ExportUserInfo struct {
	FullName string
	Email    string
}
