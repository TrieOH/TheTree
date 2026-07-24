package repos

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.TicketType, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketTypesRepo.GetByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetTicketTypeByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapTicketType(result)), nil
}
