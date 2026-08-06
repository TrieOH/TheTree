package handlers

import (
	"context"
	"reflect"
	"testing"

	"IdentityX/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// clientOnlyMethods are the Server methods that must reject project-scoped
// identities: the platform-only surface of IdentityX. Each method enforces
// models.RequireClientOnly as its first statement, so a project-scoped
// identity in the context is rejected with 401 before any request handling.
// This list is the regression net for that policy — a handler that drops
// its guard fails this test (it no longer 401s).
var clientOnlyMethods = []string{
	"ListActors", "CreateActor", "GetActor", "GetActorByEmail",
	"CreateAPIKey", "CreateCapability",
	"ListOrganizations", "CreateOrganization", "ListOrganizationMembers",
	"AddOrganizationMember", "RemoveOrganizationMember",
	"ListOrganizationProjects", "CreateOrganizationProject",
	"ListOrganizationProjectActors", "CreateOrganizationProjectActor",
	"ListOrganizationProjectMembers", "AddOrganizationProjectMember",
	"RemoveOrganizationProjectMember", "GetOrganizationProjectActor",
	"ListProjects", "CreateProject", "ListProjectMembers",
	"AddProjectMember", "RemoveProjectMember",
	"GetPlatformProfile", "UpsertPlatformProfile",
	"UpsertPlatformProfileSchema", "UpsertProjectProfileSchema",
	"ListOutdatedPlatformProfiles",
}

func projectScopedCtx() context.Context {
	pid := uuid.New()
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub: models.Subject{ID: uuid.New(), ProjectID: &pid},
	})
}

// TestClientOnlyHandlersRejectProjectScoped exercises every platform-only
// handler with a project-scoped identity and asserts it fails with 401. The
// guard runs before request handling, so a zero-value request object is
// safe; without the guard the handler would proceed (and typically error or
// panic), failing the test. The Server methods are promoted from the
// feature handler packages, so this also pins the embedded-aggregate wiring.
func TestClientOnlyHandlersRejectProjectScoped(t *testing.T) {
	srv := &Server{}
	for _, name := range clientOnlyMethods {
		t.Run(name, func(t *testing.T) {
			m := reflect.ValueOf(srv).MethodByName(name)
			if !m.IsValid() {
				t.Fatalf("Server has no method %s", name)
			}
			req := reflect.Zero(m.Type().In(1)) // nil request interface
			out := m.Call([]reflect.Value{reflect.ValueOf(projectScopedCtx()), req})
			err, _ := out[1].Interface().(error)
			if err == nil {
				t.Fatalf("expected 401 for project-scoped identity, got success")
			}
			if !fun.Is(err, fun.CodeUnauthorized) {
				t.Fatalf("want unauthorized error, got %v", err)
			}
		})
	}
}
