package profiles

import (
	"context"
	"encoding/json"
	"errors"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/jsonschema"

	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// errNoActiveSchema marks a scope without a schema (none configured or
// inactive). Profiles in that scope are stored unvalidated and have no
// version to migrate to.
var errNoActiveSchema = errors.New("no active profile schema")

// loadActiveSchema returns the active schema for the scope, or
// errNoActiveSchema when no schema exists or the existing one is inactive.
func (o *Operations) loadActiveSchema(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error) {
	ctx, span := telemetry.StartSpan(ctx, "loadActiveSchema")
	defer span.End()

	s, err := o.schemas.Get(ctx, projectID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
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
		return fun.ErrForbidden("insufficient permissions")
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
		return 0, fun.ErrValidation(err.Error())
	}
	return schema.Version, nil
}

// migrateOnDemand moves a profile to the active schema version when the
// read path finds it behind: documents that still validate bump the version
// pointer, documents that no longer validate keep their current version and
// get flagged as outdated for manual resolution. The write is skipped when
// the migration state would not change, so reads stay stable.
func (o *Operations) migrateOnDemand(ctx context.Context, profile *models.ActorProfile, projectID *uuid.UUID) (*models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "migrateOnDemand")
	defer span.End()

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

	version, outdated := schema.Version, false
	err = jsonschema.Validate(schema.Schema, profile.Profile)
	if err != nil {
		version, outdated = profile.SchemaVersion, true
	}
	if version == profile.SchemaVersion && outdated == profile.Outdated {
		return profile, nil
	}

	return o.profiles.SetMigrationState(ctx, profile.ActorID, version, outdated)
}
