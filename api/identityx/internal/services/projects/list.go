package projects

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"
)

func (o *Operations) ListProjects(ctx context.Context) ([]models.Project, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListProjects")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownProjects, err := o.projects.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedProjects, err := o.projects.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownProjects, joinedProjects...), nil
}
