package profiles

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"
)

func (r *Repo) Upsert(ctx context.Context, profile models.ActorProfile) (*models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "Upsert")
	defer span.End()

	result, err := database.Queries(ctx, r.q).UpsertActorProfile(ctx, sqlc.UpsertActorProfileParams{
		ActorID:       profile.ActorID,
		Profile:       profile.Profile,
		SchemaVersion: profile.SchemaVersion,
		Outdated:      profile.Outdated,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapActorProfile(result)), nil
}
