package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) Register(ctx context.Context, toRegister models.Actor) (*models.Actor, error) {
	ctx, span := repo.tracer.Start(ctx, "Register")
	defer span.End()
	sqlcActor, err := database.Queries(ctx, repo.q).RegisterActor(ctx, sqlc.RegisterActorParams{
		ProjectID:    toRegister.ProjectID,
		AuthMethod:   string(toRegister.AuthMethod),
		PasswordHash: toRegister.PasswordHash,
		Email:        toRegister.Email,
		Type:         string(toRegister.Type),
		Metadata:     toRegister.Metadata,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActor(sqlcActor)), nil
}
