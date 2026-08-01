package actors

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetByEmail(ctx context.Context, email string, projectID *uuid.UUID) (*models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByEmail")
	defer span.End()

	sqlcActor, err := database.Queries(ctx, repo.q).GetActorByEmail(ctx, sqlc.GetActorByEmailParams{
		Email:     &email,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActor(sqlcActor)), nil
}
