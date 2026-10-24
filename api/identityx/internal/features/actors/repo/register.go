package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *Repo) Register(ctx context.Context, toRegister models.Actor) (*models.Actor, error) {
	sqlcActor, err := database.Queries(ctx, repo.q).RegisterActor(ctx, sqlc.RegisterActorParams{
		ProjectID:    toRegister.ProjectID,
		AuthMethod:   string(toRegister.AuthMethod),
		PasswordHash: toRegister.PasswordHash,
		Email:        toRegister.Email,
		Type:         string(toRegister.Type),
		Metadata:     toRegister.Metadata,
		VerifiedAt:   toRegister.VerifiedAt,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapActor(sqlcActor)), nil
}
