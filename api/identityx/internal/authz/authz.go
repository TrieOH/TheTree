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
	orgs          ports.OrganizationRepo
	projects      ports.ProjectRepo
	platformRoles ports.PlatformRolesRepo
}

func New(orgs ports.OrganizationRepo, projects ports.ProjectRepo, platformRoles ports.PlatformRolesRepo) *Service {
	return &Service{orgs: orgs, projects: projects, platformRoles: platformRoles}
}

// CheckProject gates a project-scoped operation. An actor passes when they
// hold at least minRole on the project, either directly through a
// project_members row or by org membership: a project's organization is its
// enclosing scope, so an org role casts onto the project (org member →
// project member, org admin → project admin, org owner → project owner).
// The fallback is what makes org management work — the org owner who
// created a project holds no project_members row, only the org role.
//
// A missing project surfaces as NotFound (the resource does not exist); a
// missing role surfaces as Forbidden.
func (s *Service) CheckProject(ctx context.Context, userID, projectID uuid.UUID, minRole models.ProjectRole) error {
	project, err := s.projects.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	role, err := s.projects.GetRole(ctx, userID, projectID)
	if err == nil {
		return libauthz.Min(role, minRole)
	}
	if !fun.Is(err, fun.CodeNotFound) {
		return err
	}

	if project.OrganizationID != nil {
		orgRole, oerr := s.orgs.GetRole(ctx, userID, *project.OrganizationID)
		if oerr != nil {
			return libauthz.ForbiddenIfNotFound(oerr)
		}
		return libauthz.Min(models.ProjectRole(orgRole), minRole)
	}

	return libauthz.ForbiddenIfNotFound(err)
}

func (s *Service) CheckOrg(ctx context.Context, userID, orgID uuid.UUID, minRole models.OrganizationRole) error {
	role, err := s.orgs.GetRole(ctx, userID, orgID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}

// CheckPlatform gates a platform-scoped operation: the caller must hold at
// least minRole on the platform. Actors with no platform role are Forbidden.
func (s *Service) CheckPlatform(ctx context.Context, userID uuid.UUID, minRole models.PlatformRole) error {
	role, err := s.platformRoles.GetRole(ctx, userID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}
