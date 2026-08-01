package projects

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) ListMembers(ctx context.Context, projectID uuid.UUID) (members []models.ProjectMember, err error) {
	ctx, span := telemetry.StartSpan(ctx, "ListMembers")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckProject(ctx, ident.Sub.ID, projectID, nil, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = o.projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
