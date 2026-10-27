package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) Revoke(ctx context.Context, editionID, userID uuid.UUID, origin models.BadgeEmissionOrigin, reason string) error {
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
