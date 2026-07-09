package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"
)

func (repo *repo) CreateEvent(ctx context.Context, toCreate *contracts.Event) (*contracts.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.Create")
	defer span.End()

	event, err := database.Queries(ctx, repo.q).CreateEvent(ctx, sqlc.CreateEventParams{
		ID:             toCreate.ID,
		OwnerID:        toCreate.OwnerID,
		OrganizationID: toCreate.OrganizationID,
		Name:           toCreate.Name,
		Acronym:        toCreate.Acronym,
		Slug:           toCreate.Slug,
		Tagline:        toCreate.Tagline,
		Description:    toCreate.Description,
		IsSeries:       toCreate.IsSeries,
		LogoUrl:        toCreate.LogoUrl,
		ContactEmail:   toCreate.ContactEmail,
		SocialLinks:    toCreate.SocialLinks,
		CreatedBy:      toCreate.CreatedBy,
		GoauthScopeID:  toCreate.GoauthScopeID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return mapEventFromDB(&event), nil
}
