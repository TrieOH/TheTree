package repo

import (
	"context"
	"time"

	"IdentityX/internal/sqlc"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) SetVerifiedAt(ctx context.Context, actorID uuid.UUID, at time.Time) error {
	verifiedAt := at
	err := database.Queries(ctx, repo.q).UpdateActorVerifiedAt(ctx, sqlc.UpdateActorVerifiedAtParams{
		ActorID:    actorID,
		VerifiedAt: &verifiedAt,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}

func (repo *Repo) UpdatePasswordHash(ctx context.Context, actorID uuid.UUID, hash string) error {
	passwordHash := hash
	err := database.Queries(ctx, repo.q).UpdateActorPasswordHash(ctx, sqlc.UpdateActorPasswordHashParams{
		ActorID:      actorID,
		PasswordHash: &passwordHash,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
