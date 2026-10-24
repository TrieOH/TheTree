package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate *models.Event) (*models.Event, error) {
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
