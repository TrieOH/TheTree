package badges

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

var _ ports.BadgeTemplateRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("badge_template"),
	}
}

func mapBadgeTemplate(src sqlc.BadgeTemplate) models.BadgeTemplate {
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
