package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.TicketType, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketTypesRepo.ListByEdition")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListTicketTypesByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapTicketType), nil
}
