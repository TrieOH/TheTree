package queries

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Queries) ListOrgProjects(ctx context.Context, orgID uuid.UUID) ([]models.Project, error) {
	ctx, span := s.tracer.Start(ctx, "OrganizationService.ListOrgProjects")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := s.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != org.OwnerID {
		_, err = s.orgs.GetMember(ctx, ident.Sub.ID, orgID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	return s.projects.ListByOrganization(ctx, orgID)
}
