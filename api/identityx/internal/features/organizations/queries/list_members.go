package queries

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (s *Queries) ListMembers(ctx context.Context, orgID uuid.UUID) (members []models.OrganizationMember, err error) {
	ctx, span := s.tracer.Start(ctx, "OrganizationService.GetMembers")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var org *models.Organization
	org, err = s.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != org.OwnerID {
		_, err = s.orgs.GetMember(ctx, ident.Sub.ID, orgID)
		if err != nil {
			return nil, err
		}
	}

	members, err = s.orgs.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
