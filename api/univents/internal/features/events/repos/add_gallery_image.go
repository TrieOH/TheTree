package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) AddGalleryImage(ctx context.Context, id uuid.UUID, url string) (*contracts.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.AddGalleryImage")
	defer span.End()

	sqlcEvent, err := database.Queries(ctx, repo.q).AddEventGalleryImage(ctx, sqlc.AddEventGalleryImageParams{
		ID:  id,
		Url: url,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return mapEventFromDB(&sqlcEvent), nil
}
