package organizations

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) ListOrgProjects(ctx context.Context, orgID uuid.UUID) ([]models.Project, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOrgProjects")
	defer span.End()
	err := o.authz.CheckOrg(ctx, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.projects.ListByOrganization(ctx, orgID)
}
