package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByEmail(ctx context.Context, email string, projectID *uuid.UUID) (*models.Actor, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByEmail")
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
