package service

import (
	"IdentityX/models"
	"context"
)

func (o *Operations) ListProjects(ctx context.Context) ([]models.Project, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownProjects, err := o.Projects.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedProjects, err := o.Projects.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownProjects, joinedProjects...), nil
}
