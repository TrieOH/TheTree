package organizations

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) ListOrgProjectMembers(ctx context.Context, orgID, projectID uuid.UUID) ([]models.ProjectMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOrgProjectMembers")
	defer span.End()
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckProject(ctx, ident.Sub.ID, projectID, &orgID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	members, err := o.projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	orgMembers, err := o.orgs.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	for _, m := range orgMembers {
		members = append(members, models.ProjectMember{
			ActorID:   m.ActorID,
			ProjectID: projectID,
			Role:      models.ProjectRole(m.Role),
		})
	}

	return members, nil
}
