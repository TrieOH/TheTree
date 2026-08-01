package authz

import (
	"IdentityX/models"
	"IdentityX/ports"
	"context"

	libauthz "lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type Service struct {
	orgs     ports.OrganizationRepo
	projects ports.ProjectRepo
}

func New(orgs ports.OrganizationRepo, projects ports.ProjectRepo) *Service {
	return &Service{orgs: orgs, projects: projects}
}

func (s *Service) CheckProject(ctx context.Context, userID, projectID uuid.UUID, orgID *uuid.UUID, minRole models.ProjectRole) error {
	// TODO: investigate collapsing the org-scope into the project check —
	// the project knows its own organization, so the orgID param and the
	// four org-scoped call sites could go the way of informd's CheckForm.
	if orgID != nil {
		project, err := s.projects.GetByID(ctx, projectID)
		if err != nil {
			return err
		}
		if project.OrganizationID == nil || *project.OrganizationID != *orgID {
			return fun.ErrForbidden("insufficient permissions")
		}
		err = s.CheckOrg(ctx, userID, *orgID, models.OrganizationRoleMember)
		if err != nil {
			return err
		}
	}

	role, err := s.projects.GetRole(ctx, userID, projectID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}

func (s *Service) CheckOrg(ctx context.Context, userID, orgID uuid.UUID, minRole models.OrganizationRole) error {
	role, err := s.orgs.GetRole(ctx, userID, orgID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}
