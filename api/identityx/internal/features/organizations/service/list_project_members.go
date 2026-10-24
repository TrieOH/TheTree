package service

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) ListOrgProjectMembers(ctx context.Context, orgID, projectID uuid.UUID) ([]models.ProjectMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOrgProjectMembers")
	defer span.End()
	err := o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	members, err := o.Projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	orgMembers, err := o.Orgs.ListMembers(ctx, orgID)
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
