package service

import (
	"context"
	"encoding/json"
	"errors"

	"IdentityX/errs"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
	"lib/jsonschema"

	"github.com/google/uuid"
)

// errNoActiveSchema marks a scope without a schema (none configured or
// inactive). Profiles in that scope are stored unvalidated and have no
// version to migrate to.
var errNoActiveSchema = errors.New("no active profile schema")

// loadActiveSchema returns the active schema for the scope, or
// errNoActiveSchema when no schema exists or the existing one is inactive.
func (o *Operations) loadActiveSchema(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error) {
	s, err := o.Schemas.Get(ctx, projectID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errNoActiveSchema
		}
		return nil, err
	}
	if !s.Active {
		return nil, errNoActiveSchema
	}
	return s, nil
}

// requireActorInProject denies access to actors outside the project scope:
// project-scoped actors must belong to the given project, platform actors
// have no project and are never reachable through project routes.
func requireActorInProject(ctx context.Context, actors ports.ActorRepo, actorID, projectID uuid.UUID) error {
	actor, err := actors.GetByID(ctx, actorID)
	if err != nil {
		return err
	}
	if actor.ProjectID == nil || *actor.ProjectID != projectID {
		return errs.ErrActorOutsideProject
	}
	return nil
}

// validateAndStamp validates the profile against the active schema and
// returns the schema version to stamp on it (1 when no schema applies, so
// the profile is stored unvalidated). A failed validation is a client error.
func validateAndStamp(schema *models.ProjectProfileSchema, profile json.RawMessage) (int, error) {
	if schema == nil {
		return 1, nil
	}
	err := jsonschema.Validate(schema.Schema, profile)
	if err != nil {
		return 0, &errs.ProfileSchemaViolationError{Err: err}
	}
	return schema.Version, nil
}

// migrateOnDemand moves a profile to the active schema version when the
// read path finds it behind. Documents that already validate bump the
// version pointer. Documents that only carry fields the new schema forbids
// (additionalProperties: false) are auto-migrated on read: the forbidden
// fields are dropped, the pruned document is persisted at the new version.
// Only when the pruned document still does not validate — e.g. the new
// schema requires a field the old document cannot provide — does the
// profile keep its version and get flagged outdated for admin resolution.
// The write is skipped when the migration state would not change, so reads
// stay stable.
func (o *Operations) migrateOnDemand(ctx context.Context, profile *models.ActorProfile, projectID *uuid.UUID) (*models.ActorProfile, error) {
	schema, err := o.loadActiveSchema(ctx, projectID)
	if err != nil {
		if errors.Is(err, errNoActiveSchema) {
			return profile, nil
		}
		return nil, err
	}
	if profile.SchemaVersion >= schema.Version {
		return profile, nil
	}

	// fast path: the document already validates against the new schema,
	// only the version pointer needs bumping
	err = jsonschema.Validate(schema.Schema, profile.Profile)
	if err == nil {
		return o.Profiles.SetMigrationState(ctx, profile.ActorID, schema.Version, false)
	}

	// auto-migrate: drop the fields the new schema forbids, re-validate,
	// and persist the pruned document at the new version
	pruned, perr := pruneToSchema(schema.Schema, profile.Profile)
	if perr == nil {
		err := jsonschema.Validate(schema.Schema, pruned)
		if err == nil {
			return o.Profiles.Upsert(ctx, models.ActorProfile{
				ActorID:       profile.ActorID,
				Handle:        profile.Handle,
				PfpURL:        profile.PfpURL,
				Profile:       pruned,
				SchemaVersion: schema.Version,
				Outdated:      false,
			})
		}
	}

	// cannot auto-migrate (e.g. a new required field): keep the current
	// version and flag the profile for admin resolution
	if profile.Outdated {
		return profile, nil
	}
	return o.Profiles.SetMigrationState(ctx, profile.ActorID, profile.SchemaVersion, true)
}
