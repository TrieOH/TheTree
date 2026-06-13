package queries

import (
	"IdentityX/models"
	"context"
)

func (s *Queries) ListProjects(ctx context.Context) ([]models.Project, error) {
	ctx, span := s.tracer.Start(ctx, "ProjectService.ListProjects")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownProjects, err := s.projects.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedProjects, err := s.projects.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownProjects, joinedProjects...), nil
}
