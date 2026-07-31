package authz

import (
	"IdentityX/models"
	"IdentityX/ports"
	"context"

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
	if orgID != nil {
		err := s.checkProjectOrgAccess(ctx, userID, projectID, *orgID)
		if err != nil {
			return err
		}
	}

	return s.checkProjectAccess(ctx, userID, projectID, minRole)
}

func (s *Service) checkProjectOrgAccess(ctx context.Context, userID, projectID, orgID uuid.UUID) error {
	project, err := s.projects.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project.OrganizationID == nil || *project.OrganizationID != orgID {
		return fun.ErrForbidden("insufficient permissions")
	}

	org, err := s.orgs.GetByID(ctx, orgID)
	if err != nil {
		return err
	}
	if userID == org.OwnerID {
		return nil
	}

	_, err = s.orgs.GetMember(ctx, userID, orgID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err == nil {
		return nil
	}
	return nil
}

func (s *Service) checkProjectAccess(ctx context.Context, userID, projectID uuid.UUID, minRole models.ProjectRole) error {
	project, err := s.projects.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	if userID == project.OwnerID {
		return nil
	}

	member, err := s.projects.GetMember(ctx, userID, projectID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrForbidden("insufficient permissions")
	}
	if !member.Role.Minimum(minRole) {
		return fun.ErrForbidden("insufficient permissions")
	}
	return nil
}

func (s *Service) CheckOrg(ctx context.Context, userID, orgID uuid.UUID, minRole models.OrganizationRole) error {
	org, err := s.orgs.GetByID(ctx, orgID)
	if err != nil {
		return err
	}
	if userID == org.OwnerID {
		return nil
	}

	member, err := s.orgs.GetMember(ctx, userID, orgID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if err != nil {
		return fun.ErrForbidden("insufficient permissions")
	}
	if !member.Role.Minimum(minRole) {
		return fun.ErrForbidden("insufficient permissions")
	}
	return nil
}
