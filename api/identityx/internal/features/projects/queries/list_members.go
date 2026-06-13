package queries

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (s *Queries) ListMembers(ctx context.Context, projectID uuid.UUID) (members []models.ProjectMember, err error) {
	ctx, span := s.tracer.Start(ctx, "ProjectService.GetMembers")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var org *models.Project
	org, err = s.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != org.OwnerID {
		_, err = s.projects.GetMember(ctx, ident.Sub.ID, projectID)
		if err != nil {
			return nil, err
		}
	}

	members, err = s.projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
