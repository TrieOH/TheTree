package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"
)

func (repo *repo) Create(ctx context.Context, toCreate *models.Event) (*models.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.Create")
	defer span.End()

	event, err := database.Queries(ctx, repo.q).CreateEvent(ctx, sqlc.CreateEventParams{
		OwnerID:      toCreate.OwnerID,
		FullName:     toCreate.FullName,
		Acronym:      toCreate.Acronym,
		Slug:         toCreate.Slug,
		Description:  toCreate.Description,
		Status:       string(toCreate.Status),
		ContactEmail: toCreate.ContactEmail,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapEvent(event)), nil
}
