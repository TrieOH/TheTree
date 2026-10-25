package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate *models.TicketType) (*models.TicketType, error) {
	result, err := database.Queries(ctx, repo.q).CreateTicketType(ctx, sqlc.CreateTicketTypeParams{
		EditionID:   toCreate.EditionID,
		Name:        toCreate.Name,
		Description: toCreate.Description,
		AccessLevel: toCreate.AccessLevel,
		Price:       toCreate.PriceCents,
		MaxQuantity: toCreate.MaxQuantity,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapTicketType(result)), nil
}
