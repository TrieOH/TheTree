package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"
)

func (repo *editionsRepo) Create(ctx context.Context, toCreate *contracts.Edition) (*contracts.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Create")
	defer span.End()

	edition, err := database.Queries(ctx, repo.q).CreateEdition(ctx, sqlc.CreateEditionParams{
		ID:                   toCreate.ID,
		EventID:              toCreate.EventID,
		GoauthScopeID:        toCreate.GoauthScopeID,
		Type:                 sqlc.EditionType(toCreate.Type),
		EditionName:          toCreate.EditionName,
		Tagline:              toCreate.Tagline,
		Description:          toCreate.Description,
		Status:               sqlc.EditionStatus(contracts.EditionStatusDraft),
		RegistrationOpensAt:  toCreate.RegistrationOpensAt,
		RegistrationClosesAt: toCreate.RegistrationClosesAt,
		MonetaryType:         sqlc.EditionMonetaryType(toCreate.MonetaryType),
		StartsAt:             toCreate.StartsAt,
		EndsAt:               toCreate.EndsAt,
		Timezone:             toCreate.Timezone,
		LocationName:         toCreate.LocationName,
		LocationAddress:      toCreate.LocationAddress,
		LogoUrl:              toCreate.LogoUrl,
		BannerUrl:            toCreate.BannerUrl,
		ContactEmail:         toCreate.ContactEmail,
		ContactPhone:         toCreate.ContactPhone,
		OrganizerName:        toCreate.OrganizerName,
		CreatedBy:            toCreate.CreatedBy,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return mapEditionFromDB(&edition), nil
}
