package queries

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (q *Queries) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	ctx, span := q.tracer.Start(ctx, "CapabilityService.List")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var project *models.Project
	project, err = q.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != project.OwnerID {
		_, err = q.projects.GetMember(ctx, ident.Sub.ID, projectID)
		if err != nil {
			return nil, err
		}
	}

	return q.capabilities.List(ctx, projectID)
}
