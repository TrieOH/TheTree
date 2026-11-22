package actors

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

	"github.com/ovechkin-dm/mockio/mock"
)

func newReadOps(t *testing.T) (*Operations, ports.ActorRepo) {
	t.Helper()
	mock.SetUp(t)
	actors := mock.Mock[ports.ActorRepo]()
	ops := New(Deps{
		Actors: actors,
		Authz:  authz.New(mock.Mock[ports.OrganizationRepo](), mock.Mock[ports.ProjectRepo](), mock.Mock[ports.PlatformRolesRepo]()),
	})
	return ops, actors
}

func ctxIn(projectID uuid.UUID, id uuid.UUID, typ models.ActorType, cred models.CredentialType) context.Context {
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub:  models.Subject{ID: id, ProjectID: &projectID, Type: typ},
		Cred: models.Credential{Type: cred},
	})
}

// A project's end users can't list each other; its backends can.
func TestListActorsIsForBackends(t *testing.T) {
	project := uuid.New()
	cases := []struct {
		name string
		typ  models.ActorType
		cred models.CredentialType
		want error
	}{
		{"human with a token", models.HumanActorType, models.TokenCredentialType, libauthz.ErrInsufficientRole},
		{"human's api key", models.HumanActorType, models.APIKeyCredentialType, nil},
		{"service with a token", models.ServiceActorType, models.TokenCredentialType, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ops, actors := newReadOps(t)
			mock.When(actors.List(mock.AnyContext(), mock.Equal(project))).ThenReturn([]models.Actor{}, nil)

			_, err := ops.List(ctxIn(project, uuid.New(), tc.typ, tc.cred), project)
			if !errors.Is(err, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, err)
			}
		})
	}
}

func TestEndUserReadsOnlyThemselves(t *testing.T) {
	project := uuid.New()
	ops, actors := newReadOps(t)
	self := &models.Actor{ID: uuid.New(), ProjectID: &project}
	mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(self.ID))).ThenReturn(self, nil)
	ctx := ctxIn(project, self.ID, models.HumanActorType, models.TokenCredentialType)

	got, err := ops.GetByID(ctx, self.ID, project)
	if err != nil || got.ID != self.ID {
		t.Fatalf("self read: got %v, %v", got, err)
	}
	_, err = ops.GetByID(ctx, uuid.New(), project)
	if !errors.Is(err, libauthz.ErrInsufficientRole) {
		t.Fatalf("other read: expected insufficient role, got %v", err)
	}
}

// Platform actors have no project, so no project route reaches them.
func TestGetByIDRefusesPlatformActors(t *testing.T) {
	project := uuid.New()
	ops, actors := newReadOps(t)
	platform := &models.Actor{ID: uuid.New()}
	mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(platform.ID))).ThenReturn(platform, nil)

	_, err := ops.GetByID(ctxIn(project, uuid.New(), models.ServiceActorType, models.APIKeyCredentialType), platform.ID, project)
	if !errors.Is(err, errs.ErrActorOutsideProject) {
		t.Fatalf("expected actor outside project, got %v", err)
	}
}
