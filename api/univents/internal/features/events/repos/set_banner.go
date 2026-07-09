package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) SetBanner(ctx context.Context, id uuid.UUID, url string) (*contracts.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.SetBanner")
	defer span.End()

	sqlcEvent, err := database.Queries(ctx, repo.q).SetEventBanner(ctx, sqlc.SetEventBannerParams{
		ID:  id,
		Url: url,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return mapEventFromDB(&sqlcEvent), nil
}
