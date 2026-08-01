package ticket_types

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

var _ ports.TicketTypeRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("ticket_type"),
	}
}

func mapTicketType(src sqlc.TicketType) models.TicketType {
	return models.TicketType{
		ID:          src.ID,
		EditionID:   src.EditionID,
		Name:        src.Name,
		Description: src.Description,
		AccessLevel: src.AccessLevel,
		PriceCents:  src.Price,
		MaxQuantity: src.MaxQuantity,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}
