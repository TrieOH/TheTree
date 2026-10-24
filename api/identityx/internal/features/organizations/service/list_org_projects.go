package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) ListOrgProjects(ctx context.Context, orgID uuid.UUID) ([]models.Project, error) {
	err := o.Authz.CheckOrg(ctx, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Projects.ListByOrganization(ctx, orgID)
}
