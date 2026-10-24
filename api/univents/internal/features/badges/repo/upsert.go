package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) Upsert(ctx context.Context, emission *models.BadgeEmission) (*models.BadgeEmission, error) {
	row, err := database.Queries(ctx, repo.q).UpsertBadgeEmission(ctx, sqlc.UpsertBadgeEmissionParams{
		EditionID:      emission.EditionID,
		UserID:         emission.UserID,
		Origin:         string(emission.Origin),
		RegistrationID: emission.RegistrationID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapBadgeEmission(row)), nil
}
