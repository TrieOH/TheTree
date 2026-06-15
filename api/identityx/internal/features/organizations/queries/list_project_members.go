package queries

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Queries) ListOrgProjectMembers(ctx context.Context, orgID, projectID uuid.UUID) ([]models.ProjectMember, error) {
	ctx, span := s.tracer.Start(ctx, "OrganizationService.ListOrgProjectMembers")
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
			_, err = s.projects.GetMember(ctx, ident.Sub.ID, projectID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return nil, err
			}
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	members, err := s.projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	orgMembers, err := s.orgs.ListMembers(ctx, orgID)
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
