package ticket_types

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate *models.TicketType) (*models.TicketType, error) {
	ctx, span := telemetry.StartSpan(ctx, "TicketTypesRepo.Create")
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
