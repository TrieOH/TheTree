package commands

import (
	"IdentityX/models"
	"context"
	"lib/jsonschema"

	"lib/telemetry"

	"github.com/MintzyG/fun"
)

// UpsertPlatformProfile upserts a platform-scoped actor's profile (project_id is NULL).
// Only the actor themselves can edit; validates against the platform schema.
func (c *Commands) UpsertPlatformProfile(ctx context.Context, payload models.UpsertProfileInput) (*models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpsertPlatformProfile")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != payload.ActorID {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	// only platform schema applies
	s, err := c.schemas.Get(ctx, nil)
	if err == nil && s.Active {
		err := jsonschema.Validate(s.Schema, payload.Profile)
		if err != nil {
			return nil, fun.ErrValidation(err.Error())
		}
	}

	return c.profiles.Upsert(ctx, models.ActorProfile{
		ActorID: payload.ActorID,
		Profile: payload.Profile,
	})
}
