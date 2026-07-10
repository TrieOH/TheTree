package commands

import (
	"IdentityX/models"
	"context"
	"lib/jsonschema"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) UpsertProfile(ctx context.Context, payload models.UpsertProfileInput, projectID uuid.UUID) (*models.ActorProfile, error) {
	ctx, span := c.tracer.Start(ctx, "UpsertProfile")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	project, err := c.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// self or project admin
	if ident.Sub.ID != project.OwnerID && ident.Sub.ID != payload.ActorID {
		member, err := c.projects.GetMember(ctx, ident.Sub.ID, projectID)
		if err != nil {
			return nil, err
		}
		if member.Role != models.ProjectRoleAdmin {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	// load project schema, fall back to platform schema
	schema, err := c.loadActiveSchema(ctx, &projectID)
	if err != nil {
		return nil, err
	}

	if schema != nil {
		if err := jsonschema.Validate(schema, payload.Profile); err != nil {
			return nil, fun.ErrValidation(err.Error())
		}
	}

	return c.profiles.Upsert(ctx, models.ActorProfile{
		ActorID: payload.ActorID,
		Profile: payload.Profile,
	})
}

// loadActiveSchema returns the project's active schema, falling back to platform schema.
// Returns nil (no schema) if neither exists or is inactive.
func (c *Commands) loadActiveSchema(ctx context.Context, projectID *uuid.UUID) ([]byte, error) {
	s, err := c.schemas.Get(ctx, projectID)
	if err == nil && s.Active {
		return s.Schema, nil
	}

	// platform fallback
	s, err = c.schemas.Get(ctx, nil)
	if err != nil {
		return nil, nil // no platform schema, passthrough
	}
	if !s.Active {
		return nil, nil
	}
	return s.Schema, nil
}
