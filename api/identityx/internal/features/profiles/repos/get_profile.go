package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
	"lib/telemetry"
)

func (r *Repo) Get(ctx context.Context, actorID uuid.UUID) (*models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "Get")
	defer span.End()

	result, err := database.Queries(ctx, r.q).GetActorProfile(ctx, actorID)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapActorProfile(result)), nil
}
