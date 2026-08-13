package actors

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) List(ctx context.Context, projectID uuid.UUID) ([]models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "List")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	if err := o.requireActorReadAccess(ctx, ident, projectID); err != nil {
		return nil, err
	}

	return o.actors.List(ctx, projectID)
}
