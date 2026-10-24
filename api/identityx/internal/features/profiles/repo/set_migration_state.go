package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

// SetMigrationState bumps the schema version pointer or toggles the
// outdated flag without touching the profile document itself.
func (r *Repo) SetMigrationState(ctx context.Context, actorID uuid.UUID, schemaVersion int, outdated bool) (*models.ActorProfile, error) {
	result, err := database.Queries(ctx, r.q).SetActorProfileMigrationState(ctx, sqlc.SetActorProfileMigrationStateParams{
		ActorID:       actorID,
		SchemaVersion: schemaVersion,
		Outdated:      outdated,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapActorProfile(result)), nil
}
