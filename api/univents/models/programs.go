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

type CreateProgramRequest struct {
	Kind           ProgramKind `json:"kind"             validate:"required,oneof=activity checkpoint"`
	Name           string      `json:"name"             validate:"required,min=2"`
	Description    *string     `json:"description"`
	MinAccessLevel *int        `json:"min_access_level" validate:"omitempty,gte=0"`
	StaffOnly      bool        `json:"staff_only"`
	Price          *int64      `json:"price"`
}

func (r CreateProgramRequest) ToInput(editionID uuid.UUID) CreateProgramInput {
	return CreateProgramInput{
		EditionID:      editionID,
		Kind:           r.Kind,
		Name:           r.Name,
		Description:    r.Description,
		MinAccessLevel: r.MinAccessLevel,
		StaffOnly:      r.StaffOnly,
		Price:          r.Price,
	}
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

type PatchProgramRequest struct {
	Kind           ProgramKind `json:"kind"             validate:"required,oneof=activity checkpoint"`
	Name           string      `json:"name"             validate:"required,min=2"`
	Description    *string     `json:"description"`
	MinAccessLevel *int        `json:"min_access_level" validate:"omitempty,gte=0"`
	StaffOnly      bool        `json:"staff_only"`
	Price          *int64      `json:"price"`
	BannerURL      *string     `json:"banner_url"`
}

func (r PatchProgramRequest) ToInput(programID uuid.UUID) PatchProgramInput {
	return PatchProgramInput{
		ProgramID:      programID,
		Kind:           r.Kind,
		Name:           r.Name,
		Description:    r.Description,
		MinAccessLevel: r.MinAccessLevel,
		StaffOnly:      r.StaffOnly,
		Price:          r.Price,
		BannerURL:      r.BannerURL,
	}
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

type CreateProgramOccurrenceRequest struct {
	StartsAt    time.Time `json:"starts_at"    validate:"required"`
	EndsAt      time.Time `json:"ends_at"      validate:"required"`
	MaxCapacity *int      `json:"max_capacity" validate:"omitempty,gt=0"`
}

func (r CreateProgramOccurrenceRequest) ToInput(programID uuid.UUID) CreateProgramOccurrenceInput {
	return CreateProgramOccurrenceInput{
		ProgramID:   programID,
		StartsAt:    r.StartsAt,
		EndsAt:      r.EndsAt,
		MaxCapacity: r.MaxCapacity,
	}
}

type CreateProgramOccurrenceInput struct {
	ProgramID   uuid.UUID
	StartsAt    time.Time
	EndsAt      time.Time
	MaxCapacity *int
}

type PatchProgramOccurrenceRequest struct {
	StartsAt    time.Time `json:"starts_at"    validate:"required"`
	EndsAt      time.Time `json:"ends_at"      validate:"required"`
	MaxCapacity *int      `json:"max_capacity" validate:"omitempty,gt=0"`
}

func (r PatchProgramOccurrenceRequest) ToInput(occurrenceID uuid.UUID) PatchProgramOccurrenceInput {
	return PatchProgramOccurrenceInput{
		OccurrenceID: occurrenceID,
		StartsAt:     r.StartsAt,
		EndsAt:       r.EndsAt,
		MaxCapacity:  r.MaxCapacity,
	}
}

type PatchProgramOccurrenceInput struct {
	OccurrenceID uuid.UUID
	StartsAt     time.Time
	EndsAt       time.Time
	MaxCapacity  *int
}
