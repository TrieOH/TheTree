package programs

import (
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ProgramParticipationRepo = (*Repo)(nil)

var (
	_ ports.ProgramRepo           = (*Repo)(nil)
	_ ports.ProgramOccurrenceRepo = (*Repo)(nil)
)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("program"),
	}
}

func mapProgram(src sqlc.Program) models.Program {
	return models.Program{
		ID:             src.ID,
		EditionID:      src.EditionID,
		Kind:           models.ProgramKind(src.Kind),
		Name:           src.Name,
		Description:    src.Description,
		MinAccessLevel: src.MinAccessLevel,
		StaffOnly:      src.StaffOnly,
		Price:          &src.Price,
		BannerURL:      src.BannerUrl,
		CreatedAt:      src.CreatedAt,
		UpdatedAt:      src.UpdatedAt,
		DeletedAt:      src.DeletedAt,
	}
}

func mapProgramOccurrence(src sqlc.ProgramOccurrence) models.ProgramOccurrence {
	return models.ProgramOccurrence{
		ID:          src.ID,
		ProgramID:   src.ProgramID,
		EditionID:   src.EditionID,
		StartsAt:    src.StartsAt,
		EndsAt:      src.EndsAt,
		MaxCapacity: src.MaxCapacity,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}

func mapParticipation(src sqlc.ProgramParticipation) models.ProgramParticipation {
	return models.ProgramParticipation{
		ID:             src.ID,
		EditionID:      src.EditionID,
		OccurrenceID:   src.OccurrenceID,
		RegistrationID: src.RegistrationID,
		Status:         models.ProgramParticipationStatus(src.Status),
		CreatedAt:      src.CreatedAt,
		UpdatedAt:      src.UpdatedAt,
	}
}

func mapMyParticipation(src sqlc.ListActiveProgramParticipationsByEditionAndRegistrationRow) models.MyParticipation {
	return models.MyParticipation{
		ID:             src.ParticipationID,
		EditionID:      src.EditionID,
		OccurrenceID:   src.OccurrenceID,
		RegistrationID: src.RegistrationID,
		Status:         models.ProgramParticipationStatus(src.Status),
		CreatedAt:      src.ParticipationCreatedAt,
		UpdatedAt:      src.ParticipationUpdatedAt,
		Program: models.Program{
			ID:             src.ProgramID,
			EditionID:      src.EditionID,
			Kind:           models.ProgramKind(src.ProgramKind),
			Name:           src.ProgramName,
			Description:    src.ProgramDescription,
			MinAccessLevel: src.ProgramMinAccessLevel,
			StaffOnly:      src.ProgramStaffOnly,
			Price:          &src.ProgramPrice,
			BannerURL:      src.ProgramBannerUrl,
		},
		Occurrence: models.ProgramOccurrence{
			ID:          src.OccurrenceID,
			ProgramID:   src.ProgramID,
			EditionID:   src.EditionID,
			StartsAt:    src.OccurrenceStartsAt,
			EndsAt:      src.OccurrenceEndsAt,
			MaxCapacity: src.OccurrenceMaxCapacity,
		},
	}
}

func mapParticipant(src sqlc.ListProgramParticipationsByOccurrenceRow) models.ProgramParticipant {
	return models.ProgramParticipant{
		ID:             src.ParticipationID,
		OccurrenceID:   src.OccurrenceID,
		RegistrationID: src.RegistrationID,
		Status:         models.ProgramParticipationStatus(src.Status),
		CreatedAt:      src.ParticipationCreatedAt,
		AttendeeUserID: src.AttendeeUserID,
		AttendeeEmail:  src.AttendeeEmail,
		AttendeeName:   src.AttendeeName,
	}
}

func priceValue(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
