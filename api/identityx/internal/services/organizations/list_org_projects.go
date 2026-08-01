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
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckOrg(ctx, ident.Sub.ID, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.projects.ListByOrganization(ctx, orgID)
}
