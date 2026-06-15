package queries

import (
	"IdentityX/models"
	"context"
)

func (s *Queries) ListOrgs(ctx context.Context) ([]models.Organization, error) {
	ctx, span := s.tracer.Start(ctx, "OrganizationService.ListOrgs")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownOrgs, err := s.orgs.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedOrgs, err := s.orgs.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownOrgs, joinedOrgs...), nil
}
