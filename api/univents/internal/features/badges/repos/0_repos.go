package repos

import (
	sqlc2 "univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"lib/database"
)

type Repo struct {
	q   *sqlc2.Queries
	dbe database.ErrorHandler
}

var _ ports.BadgeTemplateRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("badge_template"),
	}
}

func mapBadgeTemplate(src sqlc2.BadgeTemplate) models.BadgeTemplate {
	return models.BadgeTemplate{
		ID:           src.ID,
		EditionID:    src.EditionID,
		TicketTypeID: src.TicketTypeID,
		Name:         src.Name,
		DesignData:   src.DesignData,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
		DeletedAt:    src.DeletedAt,
	}
}
