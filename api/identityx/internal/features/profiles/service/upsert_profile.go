package service

import (
	"errors"

	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

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
