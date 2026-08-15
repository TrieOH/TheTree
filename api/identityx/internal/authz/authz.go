package authz

import (
	"context"

	"IdentityX/models"
	"IdentityX/ports"

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

// callerID resolves the authenticated actor from the context identity the
// auth middleware wrote. Every check resolves the caller itself, so the
// interface carries no caller parameter — there is nothing to pass wrong,
// and a caller cannot be checked against a subject other than itself.
func callerID(ctx context.Context) (uuid.UUID, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return ident.Sub.ID, nil
}

// ── Scope checkers ───────────────────────────────────────────────────────

// requirePlatformClient is the "platform-only" scope checker: the caller
// must be an authenticated platform-level client — human, service, or
// machine with no project scoping. Project-scoped actors and unauthenticated
// contexts are rejected. It is the scope half of the Access-check module,
// enforced by the spec-derived chain, and directly testable as a predicate.
func (s *Service) requirePlatformClient(ctx context.Context) error {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return err
	}
	if ident.Sub.ProjectID != nil {
		return fun.ErrUnauthorized("platform-level authentication required")
	}
	return nil
}

// ScopeCheckers returns the registered scope checkers, keyed by the x-scope
// value the spec declares per operation. The resolver validates every
// x-scope an operation declares against this registry and fails boot on a
// miss. Adding a scope is one spec line plus one entry here — the generic
// mechanism in lib/authz needs no further wiring.
func (s *Service) ScopeCheckers() map[string]libauthz.ScopeChecker {
	return map[string]libauthz.ScopeChecker{
		"platform-only": s.requirePlatformClient,
	}
}

// CheckProject gates a project-scoped operation. The caller is resolved
// from the context identity; they pass when they hold at least minRole on
// the project, either directly through a project_members row or by org
// membership: a project's organization is its enclosing scope, so an org
// role casts onto the project (org member → project member, org admin →
// project admin, org owner → project owner). The fallback is what makes
// org management work — the org owner who created a project holds no
// project_members row, only the org role.
//
// A missing project surfaces as NotFound (the resource does not exist); a
// missing role surfaces as Forbidden.
func (s *Service) CheckProject(ctx context.Context, projectID uuid.UUID, minRole models.ProjectRole) error {
	userID, err := callerID(ctx)
	if err != nil {
		return err
	}

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

// CheckOrg gates an organization-scoped operation; the caller is resolved
// from the context identity.
func (s *Service) CheckOrg(ctx context.Context, orgID uuid.UUID, minRole models.OrganizationRole) error {
	userID, err := callerID(ctx)
	if err != nil {
		return err
	}

	role, err := s.orgs.GetRole(ctx, userID, orgID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}

// CheckPlatform gates a platform-scoped operation: the caller (resolved
// from the context identity) must hold at least minRole on the platform.
// Actors with no platform role are Forbidden.
func (s *Service) CheckPlatform(ctx context.Context, minRole models.PlatformRole) error {
	userID, err := callerID(ctx)
	if err != nil {
		return err
	}

	role, err := s.platformRoles.GetRole(ctx, userID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}
