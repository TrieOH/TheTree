package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
	"lib/telemetry"
)

func (repo *Repo) GetProjectServiceAccount(ctx context.Context, id uuid.UUID) (*models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetProjectServiceAccount")
	defer span.End()

	sqlcActor, err := database.Queries(ctx, repo.q).GetProjectServiceAccount(ctx, &id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActor(sqlcActor)), nil
}
