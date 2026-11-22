package projects

import (
	"context"
	"errors"
	"testing"
	"uuid"

	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"
	libauthz "lib/authz"
	"lib/database"

	"github.com/ovechkin-dm/mockio/mock"
)

// newMemberOps wires a project the caller holds callerRole on, and a target
// actor found by email.
func newMemberOps(t *testing.T, callerRole models.ProjectRole) (*Operations, ports.ProjectRepo, context.Context, *models.Project, *models.Actor) {
	t.Helper()
	mock.SetUp(t)
	caller := uuid.New()
	target := &models.Actor{ID: uuid.New()}
	project := &models.Project{ID: uuid.New(), OwnerID: uuid.New()}
	actors := mock.Mock[ports.ActorRepo]()
	projects := mock.Mock[ports.ProjectRepo]()
	mock.When(actors.GetByEmail(mock.AnyContext(), mock.Any[string](), mock.Any[*uuid.UUID]())).ThenReturn(target, nil)
	mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(target.ID))).ThenReturn(target, nil)
	mock.When(projects.GetByID(mock.AnyContext(), mock.Equal(project.ID))).ThenReturn(project, nil)
	mock.When(projects.GetRole(mock.AnyContext(), mock.Equal(caller), mock.Equal(project.ID))).ThenReturn(callerRole, nil)
	ops := New(Deps{
		Projects: projects,
		Actors:   actors,
		Authz:    authz.New(mock.Mock[ports.OrganizationRepo](), projects, mock.Mock[ports.PlatformRolesRepo]()),
	})
	ctx := models.WithIdentity(context.Background(), &models.Identity{
		Sub: models.Subject{ID: caller, Type: models.HumanActorType},
	})
	return ops, projects, ctx, project, target
}

func TestProjectAdminCannotGrantOwner(t *testing.T) {
	ops, projects, ctx, project, target := newMemberOps(t, models.ProjectRoleAdmin)
	mock.When(projects.GetMember(mock.AnyContext(), mock.Equal(target.ID), mock.Equal(project.ID))).
		ThenReturn(nil, database.ErrNotFound)

	err := ops.AddMember(ctx, models.AddProjectMemberInput{
		ProjectID: project.ID, ActorEmail: "x@y.z", Role: models.ProjectRoleOwner,
	})
	if !errors.Is(err, libauthz.ErrInsufficientRole) {
		t.Fatalf("expected insufficient role, got %v", err)
	}
}

func TestProjectAdminCannotRemoveAnAdmin(t *testing.T) {
	ops, projects, ctx, project, target := newMemberOps(t, models.ProjectRoleAdmin)
	mock.When(projects.GetMember(mock.AnyContext(), mock.Equal(target.ID), mock.Equal(project.ID))).
		ThenReturn(&models.ProjectMember{ActorID: target.ID, Role: models.ProjectRoleAdmin}, nil)

	err := ops.RemoveMember(ctx, models.RemoveProjectMemberInput{ProjectID: project.ID, ActorEmail: "x@y.z"})
	if !errors.Is(err, libauthz.ErrInsufficientRole) {
		t.Fatalf("expected insufficient role, got %v", err)
	}
}
