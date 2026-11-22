package organizations

import (
	"context"
	"errors"
	"testing"
	"uuid"

	"IdentityX/errs"
	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	libauthz "lib/authz"
	"lib/database"

	"github.com/ovechkin-dm/mockio/mock"
)

type accessFixture struct {
	caller  uuid.UUID
	target  *models.Actor
	org     *models.Organization
	project *models.Project
	orgs    ports.OrganizationRepo
	projs   ports.ProjectRepo
	ops     *Operations
}

// newAccess wires an org the caller holds callerRole in, a target actor
// found by email (an org member with targetRole, or none when empty), and a
// project that belongs to projectOrg.
func newAccess(t *testing.T, callerRole, targetRole models.OrganizationRole, projectOrg *uuid.UUID) *accessFixture {
	t.Helper()
	mock.SetUp(t)
	f := &accessFixture{
		caller: uuid.New(),
		target: &models.Actor{ID: uuid.New()},
		orgs:   mock.Mock[ports.OrganizationRepo](),
		projs:  mock.Mock[ports.ProjectRepo](),
	}
	f.org = &models.Organization{ID: uuid.New(), OwnerID: uuid.New()}
	f.project = &models.Project{ID: uuid.New(), OrganizationID: projectOrg, OwnerID: uuid.New()}
	actors := mock.Mock[ports.ActorRepo]()
	mock.When(actors.GetByEmail(mock.AnyContext(), mock.Any[string](), mock.Any[*uuid.UUID]())).ThenReturn(f.target, nil)
	mock.When(f.orgs.GetByID(mock.AnyContext(), mock.Equal(f.org.ID))).ThenReturn(f.org, nil)
	mock.When(f.orgs.GetRole(mock.AnyContext(), mock.Equal(f.caller), mock.Equal(f.org.ID))).ThenReturn(callerRole, nil)
	if targetRole == "" {
		mock.When(f.orgs.GetMember(mock.AnyContext(), mock.Equal(f.target.ID), mock.Equal(f.org.ID))).
			ThenReturn(nil, database.ErrNotFound)
	} else {
		mock.When(f.orgs.GetMember(mock.AnyContext(), mock.Equal(f.target.ID), mock.Equal(f.org.ID))).
			ThenReturn(&models.OrganizationMember{ActorID: f.target.ID, Role: targetRole}, nil)
	}
	mock.When(f.projs.GetByID(mock.AnyContext(), mock.Equal(f.project.ID))).ThenReturn(f.project, nil)
	f.ops = New(Deps{
		Projects: f.projs,
		Actors:   actors,
		Orgs:     f.orgs,
		Authz:    authz.New(f.orgs, f.projs, mock.Mock[ports.PlatformRolesRepo]()),
	})
	return f
}

func (f *accessFixture) ctx() context.Context {
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub: models.Subject{ID: f.caller, Type: models.HumanActorType},
	})
}

// An org admin reaches only the org's own projects: naming another org's
// project, or a standalone one, is not found.
func TestOrgProjectRoutesStayInsideTheOrg(t *testing.T) {
	other := uuid.New()
	for name, projectOrg := range map[string]*uuid.UUID{"another org": &other, "standalone": nil} {
		t.Run(name, func(t *testing.T) {
			f := newAccess(t, models.OrganizationRoleOwner, "", projectOrg)

			err := f.ops.AddProjectMember(f.ctx(), models.AddOrgProjectMemberInput{
				OrganizationID: f.org.ID, ProjectID: f.project.ID, ActorEmail: "x@y.z", Role: models.ProjectRoleMember,
			})
			if !errors.Is(err, errs.ErrProjectNotFound) {
				t.Fatalf("add: expected project not found, got %v", err)
			}
			err = f.ops.RemoveProjectMember(f.ctx(), models.RemoveOrgProjectMemberInput{
				OrganizationID: f.org.ID, ProjectID: f.project.ID, ActorEmail: "x@y.z",
			})
			if !errors.Is(err, errs.ErrProjectNotFound) {
				t.Fatalf("remove: expected project not found, got %v", err)
			}
			_, err = f.ops.ListOrgProjectMembers(f.ctx(), f.org.ID, f.project.ID)
			if !errors.Is(err, errs.ErrProjectNotFound) {
				t.Fatalf("list: expected project not found, got %v", err)
			}
		})
	}
}

func TestOrgAdminCannotGrantOwner(t *testing.T) {
	f := newAccess(t, models.OrganizationRoleAdmin, "", nil)
	f.project.OrganizationID = &f.org.ID

	err := f.ops.AddProjectMember(f.ctx(), models.AddOrgProjectMemberInput{
		OrganizationID: f.org.ID, ProjectID: f.project.ID, ActorEmail: "x@y.z", Role: models.ProjectRoleOwner,
	})
	if !errors.Is(err, libauthz.ErrInsufficientRole) {
		t.Fatalf("project owner grant: expected insufficient role, got %v", err)
	}
	err = f.ops.AddMember(f.ctx(), models.AddOrganizationMemberInput{
		OrganizationID: f.org.ID, ActorEmail: "x@y.z", Role: models.OrganizationRoleOwner,
	})
	if !errors.Is(err, libauthz.ErrInsufficientRole) {
		t.Fatalf("org owner grant: expected insufficient role, got %v", err)
	}
}

func TestOrgAdminCannotRemoveAnAdmin(t *testing.T) {
	f := newAccess(t, models.OrganizationRoleAdmin, models.OrganizationRoleAdmin, nil)

	err := f.ops.RemoveMember(f.ctx(), models.RemoveOrganizationMemberInput{
		OrganizationID: f.org.ID, ActorEmail: "x@y.z",
	})
	if !errors.Is(err, libauthz.ErrInsufficientRole) {
		t.Fatalf("expected insufficient role, got %v", err)
	}
}
