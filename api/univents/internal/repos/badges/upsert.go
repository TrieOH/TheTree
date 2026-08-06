package badges

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) Upsert(ctx context.Context, emission *models.BadgeEmission) (*models.BadgeEmission, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeEmissionsRepo.Upsert")
	defer span.End()
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
