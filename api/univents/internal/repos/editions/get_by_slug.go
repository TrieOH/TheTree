package editions

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetBySlug(ctx context.Context, eventID uuid.UUID, slug string) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionsRepo.GetBySlug")
	defer span.End()
	edition, err := database.Queries(ctx, repo.q).GetEditionBySlug(ctx, sqlc.GetEditionBySlugParams{
		Slug:    slug,
		EventID: eventID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEdition(edition)), nil
}
