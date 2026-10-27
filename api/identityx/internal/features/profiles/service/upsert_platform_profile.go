package service

import (
	"errors"

	"IdentityX/errs"
	"IdentityX/models"
	"context"
)

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
