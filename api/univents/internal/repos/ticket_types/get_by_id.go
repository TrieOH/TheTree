package ticket_types

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.TicketType, error) {
	ctx, span := telemetry.StartSpan(ctx, "TicketTypesRepo.GetByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetTicketTypeByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapTicketType(result)), nil
}

// GetByIDForUpdate is the row-lock variant used inside the checkout tx
// (split 7): serializes concurrent checkouts on the same ticket type before
// availability is checked — no oversell on the last unit. The lock is held
// until the enclosing tx commits/rolls back.
func (repo *Repo) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.TicketType, error) {
	ctx, span := telemetry.StartSpan(ctx, "TicketTypesRepo.GetByIDForUpdate")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetTicketTypeByIDForUpdate(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapTicketType(result)), nil
}
