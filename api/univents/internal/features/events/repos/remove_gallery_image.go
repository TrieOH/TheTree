package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) RemoveGalleryImage(ctx context.Context, id uuid.UUID, url string) (*contracts.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.RemoveGalleryImage")
	defer span.End()

	sqlcEvent, err := database.Queries(ctx, repo.q).RemoveEventGalleryImage(ctx, sqlc.RemoveEventGalleryImageParams{
		ID:  id,
		Url: url,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return mapEventFromDB(&sqlcEvent), nil
}
