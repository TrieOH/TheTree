package queries

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) ListOrgProjectMembers(ctx context.Context, orgID, projectID uuid.UUID) ([]models.ProjectMember, error) {
	ctx, span := q.tracer.Start(ctx, "OrganizationService.ListOrgProjectMembers")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := q.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != org.OwnerID {
		_, err = q.orgs.GetMember(ctx, ident.Sub.ID, orgID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			_, err = q.projects.GetMember(ctx, ident.Sub.ID, projectID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return nil, err
			}
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	members, err := q.projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	orgMembers, err := q.orgs.ListMembers(ctx, orgID)
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
