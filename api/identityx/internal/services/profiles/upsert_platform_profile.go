package profiles

import (
	"errors"

	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/MintzyG/fun"
)

// UpsertPlatformProfile upserts a platform-scoped actor's profile (project_id is NULL).
// Only the actor themselves can edit; validates against the platform schema.
func (o *Operations) UpsertPlatformProfile(ctx context.Context, payload models.UpsertProfileInput) (*models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpsertPlatformProfile")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != payload.ActorID {
		return nil, fun.ErrForbidden("insufficient permissions")
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

	return o.profiles.Upsert(ctx, models.ActorProfile{
		ActorID:       payload.ActorID,
		Profile:       payload.Profile,
		SchemaVersion: version,
		Outdated:      false,
	})
}
