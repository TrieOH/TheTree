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

var _ ports.EditionRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("edition"),
	}
}

func mapEdition(src sqlc.Edition) models.Edition {
	return models.Edition{
		ID:                  src.ID,
		EventID:             src.EventID,
		Name:                src.EditionName,
		Slug:                src.Slug,
		Tagline:             src.Tagline,
		Description:         src.Description,
		IsDraft:             src.IsDraft,
		RegistrationOpensAt: src.RegistrationOpensAt,
		StartsAt:            src.StartsAt,
		EndsAt:              src.EndsAt,
		LocationName:        src.LocationName,
		LocationAddress:     src.LocationAddress,
		LogoURL:             src.LogoUrl,
		BannerURL:           src.BannerUrl,
		ContactEmail:        src.ContactEmail,
		CreatedBy:           src.CreatedBy,
		CreatedAt:           src.CreatedAt,
		UpdatedAt:           src.UpdatedAt,
		DeletedAt:           src.DeletedAt,
	}
}
