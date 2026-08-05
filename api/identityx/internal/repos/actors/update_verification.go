package actors

import (
	"context"
	"time"

	"IdentityX/internal/sqlc"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) SetVerifiedAt(ctx context.Context, actorID uuid.UUID, at time.Time) error {
	ctx, span := telemetry.StartSpan(ctx, "SetVerifiedAt")
	defer span.End()

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
	ctx, span := telemetry.StartSpan(ctx, "UpdatePasswordHash")
	defer span.End()

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
