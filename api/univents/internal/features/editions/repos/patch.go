package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) Patch(ctx context.Context, id uuid.UUID, edition *models.Edition) (*models.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Patch")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).PatchEdition(ctx, sqlc.PatchEditionParams{
		EditionName:         edition.Name,
		Slug:                edition.Slug,
		Tagline:             edition.Tagline,
		Description:         edition.Description,
		RegistrationOpensAt: edition.RegistrationOpensAt,
		StartsAt:            edition.StartsAt,
		EndsAt:              edition.EndsAt,
		LocationName:        edition.LocationName,
		LocationAddress:     edition.LocationAddress,
		LogoUrl:             edition.LogoURL,
		BannerUrl:           edition.BannerURL,
		ContactEmail:        edition.ContactEmail,
		ID:                  id,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEdition(result)), nil
}
