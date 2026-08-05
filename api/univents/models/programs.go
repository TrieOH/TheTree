package models

import (
	"time"

	"github.com/google/uuid"
)

type ProgramKind string

const (
	ProgramKindActivity   ProgramKind = "activity"
	ProgramKindCheckpoint ProgramKind = "checkpoint"
)

type Program struct {
	ID             uuid.UUID   `json:"id"`
	EditionID      uuid.UUID   `json:"edition_id"`
	Kind           ProgramKind `json:"kind"`
	Name           string      `json:"name"`
	Description    *string     `json:"description"`
	MinAccessLevel *int        `json:"min_access_level"`
	StaffOnly      bool        `json:"staff_only"`
	Price          *int64      `json:"price"`
	BannerURL      *string     `json:"banner_url"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      *time.Time  `json:"updated_at"`
	DeletedAt      *time.Time  `json:"deleted_at"`
}

type ProgramOccurrence struct {
	ID          uuid.UUID  `json:"id"`
	ProgramID   uuid.UUID  `json:"program_id"`
	EditionID   uuid.UUID  `json:"edition_id"`
	StartsAt    time.Time  `json:"starts_at"`
	EndsAt      time.Time  `json:"ends_at"`
	MaxCapacity *int       `json:"max_capacity"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type ProgramParticipationStatus string

const (
	ProgramParticipationStatusRegistered ProgramParticipationStatus = "registered"
	ProgramParticipationStatusAttended   ProgramParticipationStatus = "attended"
	ProgramParticipationStatusNoShow     ProgramParticipationStatus = "no_show"
	ProgramParticipationStatusCancelled  ProgramParticipationStatus = "cancelled"
)

type ProgramParticipation struct {
	ID             uuid.UUID                  `json:"id"`
	EditionID      uuid.UUID                  `json:"edition_id"`
	OccurrenceID   uuid.UUID                  `json:"occurrence_id"`
	RegistrationID uuid.UUID                  `json:"registration_id"`
	Status         ProgramParticipationStatus `json:"status"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      *time.Time                 `json:"updated_at"`
}

type CreateProgramInput struct {
	EditionID      uuid.UUID
	Kind           ProgramKind
	Name           string
	Description    *string
	MinAccessLevel *int
	StaffOnly      bool
	Price          *int64
}

type PatchProgramInput struct {
	ProgramID      uuid.UUID
	Kind           ProgramKind
	Name           string
	Description    *string
	MinAccessLevel *int
	StaffOnly      bool
	Price          *int64
	BannerURL      *string
}

type CreateProgramOccurrenceInput struct {
	ProgramID   uuid.UUID
	StartsAt    time.Time
	EndsAt      time.Time
	MaxCapacity *int
}

type PatchProgramOccurrenceInput struct {
	OccurrenceID uuid.UUID
	StartsAt     time.Time
	EndsAt       time.Time
	MaxCapacity  *int
}
