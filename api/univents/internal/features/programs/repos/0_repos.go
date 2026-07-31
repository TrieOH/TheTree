package repos

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

func priceValue(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
