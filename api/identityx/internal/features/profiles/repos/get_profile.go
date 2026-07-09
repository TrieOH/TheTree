package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (r *profileRepo) Get(ctx context.Context, actorID uuid.UUID) (*models.ActorProfile, error) {
	ctx, span := database.Span(ctx, r.tracer, "GetProfile")
	defer span.End()
	result, err := database.Queries(ctx, r.q).GetActorProfile(ctx, actorID)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapActorProfile(result)), nil
}
