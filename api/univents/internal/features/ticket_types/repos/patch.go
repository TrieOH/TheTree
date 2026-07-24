package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) Patch(ctx context.Context, id uuid.UUID, ticketType *models.TicketType) (*models.TicketType, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketTypesRepo.Patch")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).PatchTicketType(ctx, sqlc.PatchTicketTypeParams{
		Name:        ticketType.Name,
		Description: ticketType.Description,
		AccessLevel: ticketType.AccessLevel,
		Price:       ticketType.PriceCents,
		MaxQuantity: ticketType.MaxQuantity,
		ID:          id,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapTicketType(result)), nil
}
