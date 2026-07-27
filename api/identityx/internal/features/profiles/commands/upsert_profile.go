package commands

import (
	"IdentityX/models"
	"context"
	"errors"
	"lib/jsonschema"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"lib/telemetry"
)

func (c *Commands) UpsertProfile(ctx context.Context, payload models.UpsertProfileInput, projectID uuid.UUID) (*models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpsertProfile")
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
		err := jsonschema.Validate(schema, payload.Profile)
		if err != nil {
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
	ctx, span := telemetry.StartSpan(ctx, "loadActiveSchema")
	defer span.End()

	s, err := c.schemas.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if s.Active {
		return s.Schema, nil
	}
	return nil, errors.New("no active profile schema")
}
