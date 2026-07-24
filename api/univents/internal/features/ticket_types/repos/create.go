package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"
)

func (repo *repo) Create(ctx context.Context, toCreate *models.TicketType) (*models.TicketType, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketTypesRepo.Create")
	defer span.End()
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
