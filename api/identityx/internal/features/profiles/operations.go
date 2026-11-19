// Package profiles is the profiles feature: its repo, handlers and
// operations. Build wires it.
package profiles

import (
	"IdentityX/errs"
	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lib/database"
	"lib/jsonschema"
	"uuid"
)

type Deps struct {
	Profiles ports.ProfileRepo
	Schemas  ports.ProfileSchemaRepo
	Actors   ports.ActorRepo
	Authz    *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// GetProfileByHandle returns the profile with the given handle. Handles are
// globally unique when present, so at most one profile can match. Public
// read, like the other profile get routes.
func (o *Operations) GetProfileByHandle(ctx context.Context, handle string) (*models.ActorProfile, error) {
	profile, err := o.Profiles.GetByHandle(ctx, handle)
	if err != nil {
		return nil, err
	}

	// migrate against the actor's scope (project or platform)
	actor, err := o.Actors.GetByID(ctx, profile.ActorID)
	if err != nil {
		return nil, err
	}

	return o.migrateOnDemand(ctx, profile, actor.ProjectID)
}

// GetPlatformProfile returns a platform-scoped actor's profile (project_id is NULL).
func (o *Operations) GetPlatformProfile(ctx context.Context, actorID uuid.UUID) (*models.ActorProfile, error) {
	profile, err := o.Profiles.Get(ctx, actorID)
	if err != nil {
		return nil, err
	}

	return o.migrateOnDemand(ctx, profile, nil)
}

func (o *Operations) GetProfile(ctx context.Context, actorID, projectID uuid.UUID) (*models.ActorProfile, error) {
	// public read: the project route serves project-scoped actors only
	err := requireActorInProject(ctx, o.Actors, actorID, projectID)
	if err != nil {
		return nil, err
	}

	profile, err := o.Profiles.Get(ctx, actorID)
	if err != nil {
		return nil, err
	}

	return o.migrateOnDemand(ctx, profile, &projectID)
}

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

// ListOutdatedProfiles lists profiles flagged as outdated (they failed to
// migrate to the active schema version). A nil projectID means the platform
// scope; the project variant requires a project admin.
func (o *Operations) ListOutdatedProfiles(ctx context.Context, projectID *uuid.UUID) ([]models.ActorProfile, error) {
	if projectID != nil {
		err := o.Authz.CheckProject(ctx, *projectID, models.ProjectRoleAdmin)
		if err != nil {
			return nil, err
		}
	}

	return o.Profiles.ListOutdated(ctx, projectID)
}

// pruneToSchema rewrites a profile document to satisfy the given schema by
// dropping every property the schema forbids: keys not listed under
// "properties" are removed wherever "additionalProperties" is false, and
// nested objects are pruned recursively. Fields the schema still allows are
// left untouched. The result is only meaningful if it re-validates against
// the schema; callers must re-run jsonschema.Validate.
func pruneToSchema(schema, doc json.RawMessage) (json.RawMessage, error) {
	var sch map[string]any
	err := json.Unmarshal(schema, &sch)
	if err != nil {
		return nil, fmt.Errorf("prune: invalid schema json: %w", err)
	}
	var inst any
	err = json.Unmarshal(doc, &inst)
	if err != nil {
		return nil, fmt.Errorf("prune: invalid instance json: %w", err)
	}
	pruned, err := pruneNode(sch, inst)
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(pruned)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// pruneNode drops keys from an instance object that its schema node
// forbids, recursing into nested objects and array items. It understands
// the JSON Schema draft 2020-12 shape: {"type": ..., "properties": {...},
// "additionalProperties": bool, "items": {...}}. Nodes without a
// "properties" map are returned untouched.
func pruneNode(sch, inst any) (any, error) {
	schemaNode, _ := sch.(map[string]any)
	obj, ok := inst.(map[string]any)
	if !ok {
		if arr, isArr := inst.([]any); isArr {
			items, _ := schemaNode["items"].(map[string]any)
			for i, el := range arr {
				pruned, err := pruneNode(items, el)
				if err != nil {
					return nil, err
				}
				arr[i] = pruned
			}
			return arr, nil
		}
		return inst, nil
	}

	props, _ := schemaNode["properties"].(map[string]any)
	if props == nil {
		return inst, nil
	}
	strict := false
	if ap, ok := schemaNode["additionalProperties"].(bool); ok && !ap {
		strict = true
	}

	out := make(map[string]any, len(obj))
	for key, val := range obj {
		sub, ok := props[key]
		if !ok {
			if strict {
				continue // drop the forbidden key
			}
			out[key] = val
			continue
		}
		pruned, err := pruneNode(sub, val)
		if err != nil {
			return nil, err
		}
		out[key] = pruned
	}
	return out, nil
}

// UpsertPlatformProfile upserts a platform-scoped actor's profile (project_id is NULL).
// Only the actor themselves can edit; validates against the platform schema.
func (o *Operations) UpsertPlatformProfile(ctx context.Context, payload models.UpsertProfileInput) (*models.ActorProfile, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != payload.ActorID {
		return nil, errs.ErrNotProfileOwner
	}

	schema, err := o.loadActiveSchema(ctx, nil)
	if err != nil {
		if !errors.Is(err, errNoActiveSchema) {
			return nil, err
		}
		schema = nil // no active schema: store unvalidated
	}

	version, err := validateAndStamp(schema, payload.Profile)
	if err != nil {
		return nil, err
	}

	return o.Profiles.Upsert(ctx, models.ActorProfile{
		ActorID:       payload.ActorID,
		Handle:        payload.Handle,
		PfpURL:        payload.PfpURL,
		Profile:       payload.Profile,
		SchemaVersion: version,
		Outdated:      false,
	})
}

func (o *Operations) UpsertProfile(ctx context.Context, payload models.UpsertProfileInput, projectID uuid.UUID) (*models.ActorProfile, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	// self or project admin
	if ident.Sub.ID != payload.ActorID {
		err = o.Authz.CheckProject(ctx, projectID, models.ProjectRoleAdmin)
		if err != nil {
			return nil, err
		}
	}

	err = requireActorInProject(ctx, o.Actors, payload.ActorID, projectID)
	if err != nil {
		return nil, err
	}

	schema, err := o.loadActiveSchema(ctx, &projectID)
	if err != nil {
		if !errors.Is(err, errNoActiveSchema) {
			return nil, err
		}
		schema = nil // no active schema: store unvalidated
	}

	version, err := validateAndStamp(schema, payload.Profile)
	if err != nil {
		return nil, err
	}

	return o.Profiles.Upsert(ctx, models.ActorProfile{
		ActorID:       payload.ActorID,
		Handle:        payload.Handle,
		PfpURL:        payload.PfpURL,
		Profile:       payload.Profile,
		SchemaVersion: version,
		Outdated:      false,
	})
}
