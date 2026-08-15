package profile_schemas

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"
)

func (o *Operations) UpsertSchema(ctx context.Context, payload models.UpsertProfileSchemaInput) (*models.ProjectProfileSchema, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpsertSchema")
	defer span.End()

	// Platform schema: a platform-scoped write, so it is gated on the
	// caller's platform role rather than a project membership.
	if payload.ProjectID == nil {
		ident, err := models.RequireIdentity(ctx)
		if err != nil {
			return nil, err
		}
		err = o.authz.CheckPlatform(ctx, ident.Sub.ID, models.PlatformRoleAdmin)
		if err != nil {
			return nil, err
		}
		return o.schemas.Upsert(ctx, models.ProjectProfileSchema{
			ProjectID: nil,
			Schema:    payload.Schema,
			Active:    payload.Active,
		})
	}

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckProject(ctx, ident.Sub.ID, *payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.schemas.Upsert(ctx, models.ProjectProfileSchema{
		ProjectID: payload.ProjectID,
		Schema:    payload.Schema,
		Active:    payload.Active,
	})
}
