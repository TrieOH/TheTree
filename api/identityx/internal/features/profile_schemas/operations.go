// Package profile_schemas is the profile schemas feature: its repo, handlers
// and operations. Build wires it.
package profile_schemas

import (
	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"uuid"
)

type Deps struct {
	Schemas  ports.ProfileSchemaRepo
	Projects ports.ProjectRepo
	Authz    *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// GetSchema returns the active profile schema for the scope. Public read:
// schemas shape public profiles, so anyone may fetch them.
func (o *Operations) GetSchema(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error) {
	return o.Schemas.Get(ctx, projectID)
}

func (o *Operations) UpsertSchema(ctx context.Context, payload models.UpsertProfileSchemaInput) (*models.ProjectProfileSchema, error) {
	// Platform schema: a platform-scoped write, so it is gated on the
	// caller's platform role rather than a project membership.
	if payload.ProjectID == nil {
		err := o.Authz.CheckPlatform(ctx, models.PlatformRoleAdmin)
		if err != nil {
			return nil, err
		}
		return o.Schemas.Upsert(ctx, models.ProjectProfileSchema{
			ProjectID: nil,
			Schema:    payload.Schema,
			Active:    payload.Active,
		})
	}

	err := o.Authz.CheckProject(ctx, *payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Schemas.Upsert(ctx, models.ProjectProfileSchema{
		ProjectID: payload.ProjectID,
		Schema:    payload.Schema,
		Active:    payload.Active,
	})
}
