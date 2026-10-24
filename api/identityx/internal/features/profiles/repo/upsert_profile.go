package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (r *Repo) Upsert(ctx context.Context, profile models.ActorProfile) (*models.ActorProfile, error) {
	result, err := database.Queries(ctx, r.q).UpsertActorProfile(ctx, sqlc.UpsertActorProfileParams{
		ActorID:       profile.ActorID,
		Handle:        profile.Handle,
		PfpUrl:        profile.PfpURL,
		Profile:       profile.Profile,
		SchemaVersion: profile.SchemaVersion,
		Outdated:      profile.Outdated,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapActorProfile(result)), nil
}
