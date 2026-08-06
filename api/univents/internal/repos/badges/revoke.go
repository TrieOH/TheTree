package badges

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) Revoke(ctx context.Context, editionID, userID uuid.UUID, origin models.BadgeEmissionOrigin, reason string) error {
	ctx, span := telemetry.StartSpan(ctx, "BadgeEmissionsRepo.Revoke")
	defer span.End()
	err := database.Queries(ctx, repo.q).RevokeBadgeEmission(ctx, sqlc.RevokeBadgeEmissionParams{
		EditionID: editionID,
		UserID:    userID,
		Origin:    string(origin),
		Reason:    new(reason),
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
