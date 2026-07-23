package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) ListEventMembers(ctx context.Context, eventID uuid.UUID) ([]models.EventMember, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.ListEventMembers")
	defer span.End()

	members, err := database.Queries(ctx, repo.q).ListEventMembers(ctx, eventID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(members, mapEventMember), nil
}
