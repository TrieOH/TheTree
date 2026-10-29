package profiles

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ProfileRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrProfileNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (r *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, r.q) }

func mapActorProfile(src sqlc.ActorProfile) models.ActorProfile {
	return models.ActorProfile{
		ActorID:       src.ActorID,
		Handle:        src.Handle,
		PfpURL:        src.PfpUrl,
		Profile:       src.Profile,
		SchemaVersion: src.SchemaVersion,
		Outdated:      src.Outdated,
		UpdatedAt:     src.UpdatedAt,
	}
}

// GetByHandle returns the profile whose handle matches. Handles are unique
// when present, so at most one profile can match.
func (r *Repo) GetByHandle(ctx context.Context, handle string) (*models.ActorProfile, error) {
	result, err := r.db(ctx).GetActorProfileByHandle(ctx, &handle)
	return database.One(result, err, r.dbe, mapActorProfile)
}

func (r *Repo) Get(ctx context.Context, actorID uuid.UUID) (*models.ActorProfile, error) {
	result, err := r.db(ctx).GetActorProfile(ctx, actorID)
	return database.One(result, err, r.dbe, mapActorProfile)
}

// ListOutdated returns profiles flagged as outdated. A nil projectID means
// the platform scope (actors with no project).
func (r *Repo) ListOutdated(ctx context.Context, projectID *uuid.UUID) ([]models.ActorProfile, error) {
	results, err := r.db(ctx).ListOutdatedProfiles(ctx, projectID)
	return database.Many(results, err, r.dbe, mapActorProfile)
}

// SetMigrationState bumps the schema version pointer or toggles the
// outdated flag without touching the profile document itself.
func (r *Repo) SetMigrationState(ctx context.Context, actorID uuid.UUID, schemaVersion int, outdated bool) (*models.ActorProfile, error) {
	result, err := r.db(ctx).SetActorProfileMigrationState(ctx, sqlc.SetActorProfileMigrationStateParams{
		ActorID:       actorID,
		SchemaVersion: schemaVersion,
		Outdated:      outdated,
	})
	return database.One(result, err, r.dbe, mapActorProfile)
}

func (r *Repo) Upsert(ctx context.Context, profile models.ActorProfile) (*models.ActorProfile, error) {
	result, err := r.db(ctx).UpsertActorProfile(ctx, sqlc.UpsertActorProfileParams{
		ActorID:       profile.ActorID,
		Handle:        profile.Handle,
		PfpUrl:        profile.PfpURL,
		Profile:       profile.Profile,
		SchemaVersion: profile.SchemaVersion,
		Outdated:      profile.Outdated,
	})
	return database.One(result, err, r.dbe, mapActorProfile)
}
