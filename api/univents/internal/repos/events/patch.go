package events

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) Patch(ctx context.Context, id uuid.UUID, event *models.Event) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.Patch")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).PatchEvent(ctx, sqlc.PatchEventParams{
		FullName:     event.FullName,
		Acronym:      event.Acronym,
		Slug:         event.Slug,
		Description:  event.Description,
		LogoUrl:      event.LogoURL,
		BannerUrl:    event.BannerURL,
		ContactEmail: event.ContactEmail,
		ID:           id,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(result)), nil
}
