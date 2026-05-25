package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (repo *repo) GetByEmail(ctx context.Context, email string, projectID *uuid.UUID) (*models.Actor, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByEmail")
	defer span.End()
	repo.log.Info("get by email data", zap.String("email", email), zap.Any("projectID", projectID))
	sqlcActor, err := database.Queries(ctx, repo.q).GetActorByEmail(ctx, sqlc.GetActorByEmailParams{
		Email:     &email,
		ProjectID: projectID,
	})
	repo.log.Info("get by email", zap.Error(err), zap.Any("actor", sqlcActor))
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActor(sqlcActor)), nil
}
