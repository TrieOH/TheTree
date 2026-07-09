package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (r *profileRepo) Upsert(ctx context.Context, profile models.ActorProfile) (*models.ActorProfile, error) {
	ctx, span := database.Span(ctx, r.tracer, "UpsertProfile")
	defer span.End()
	result, err := database.Queries(ctx, r.q).UpsertActorProfile(ctx, sqlc.UpsertActorProfileParams{
		ActorID: profile.ActorID,
		Profile: profile.Profile,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapActorProfile(result)), nil
}
