package profile_schemas

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) GetSchema(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetSchema")
	defer span.End()

	// platform schema: any authenticated actor can read it
	if projectID == nil {
		_, err := models.RequireIdentity(ctx)
		if err != nil {
			return nil, err
		}
		return o.schemas.Get(ctx, nil)
	}

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	// any actor scoped to the project can read its schema
	if ident.Sub.ProjectID == nil || *ident.Sub.ProjectID != *projectID {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	return o.schemas.Get(ctx, projectID)
}
