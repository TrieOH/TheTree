package api_keys

import (
	"context"
	"testing"

	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

// testOps bundles the per-test mockio mocks backing the operations.
type testOps struct {
	apiKeys  ports.APIKeysRepo
	projects ports.ProjectRepo
	orgs     ports.OrganizationRepo
	ops      *Operations
}

// newTestOps wires the operations over fresh mocks, with the caller's
// project role resolved by authz. An empty role means "no membership":
// both the project role and the org fallback resolve to not-found.
func newTestOps(t *testing.T, project *models.Project, role models.ProjectRole) *testOps {
	t.Helper()
	mock.SetUp(t)
	o := &testOps{
		apiKeys:  mock.Mock[ports.APIKeysRepo](),
		projects: mock.Mock[ports.ProjectRepo](),
		orgs:     mock.Mock[ports.OrganizationRepo](),
	}
	mock.When(o.projects.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(project, nil)
	if role == "" {
		mock.When(o.projects.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
			ThenReturn(models.ProjectRole(""), fun.Err("project member not found").NotFound())
		mock.When(o.orgs.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
			ThenReturn(models.OrganizationRole(""), fun.Err("org member not found").NotFound())
	} else {
		mock.When(o.projects.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
			ThenReturn(role, nil)
	}
	o.ops = NewOperations(
		[]byte("test-hmac-secret"),
		mock.Mock[ports.ActorRepo](),
		o.apiKeys,
		mock.Mock[ports.CapabilityRepo](),
		o.projects,
		authz.New(o.orgs, o.projects, mock.Mock[ports.PlatformRolesRepo]()),
	)
	return o
}

// ctxAs seeds the context identity the auth middleware would have written,
// so the Access-check module resolves the caller the way production does.
func ctxAs(actorID uuid.UUID) context.Context {
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub: models.Subject{ID: actorID, Type: models.HumanActorType},
	})
}

func projectFixture() *models.Project {
	return &models.Project{ID: uuid.New(), BrandSlug: "testco", Name: "Test Co"}
}

func baseInput(projectID uuid.UUID) models.CreateAPIKeyInput {
	return models.CreateAPIKeyInput{
		Name:      "deploy key",
		Env:       "prod",
		ProjectID: &projectID,
	}
}

func withSubject(in models.CreateAPIKeyInput, subjectID uuid.UUID) models.CreateAPIKeyInput {
	in.SubjectID = &subjectID
	return in
}

// stubCreate answers APIKeysRepo.Create by assigning an ID, the way the
// sql adapter would.
func stubCreate(o *testOps) {
	mock.When(o.apiKeys.Create(mock.AnyContext(), mock.Any[models.APIKey]())).
		ThenAnswer(func(args []any) []any {
			k := args[1].(models.APIKey)
			k.ID = uuid.New()
			return []any{&k, nil}
		})
}

// TestCreateAPIKeySelfBoundMintRequiresAdmin pins the fix: a platform-level
// client minting a self-bound key in a project they hold no role in is
// forbidden. Before the fix the membership check only fired for
// impersonation mints, so the self-bound path silently bypassed the seam.
func TestCreateAPIKeySelfBoundMintRequiresAdmin(t *testing.T) {
	project := projectFixture() // no org: no fallback
	o := newTestOps(t, project, "")

	_, _, err := o.ops.Create(ctxAs(uuid.New()), baseInput(project.ID))
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("expected forbidden for non-member self-bound mint, got %v", err)
	}
}

func TestCreateAPIKeySelfBoundAdminPasses(t *testing.T) {
	project := projectFixture()
	o := newTestOps(t, project, models.ProjectRoleAdmin)
	stubCreate(o)
	caller := uuid.New()

	key, raw, err := o.ops.Create(ctxAs(caller), baseInput(project.ID))
	if err != nil {
		t.Fatalf("expected admin self-bound mint to pass, got %v", err)
	}
	if key == nil || raw == "" {
		t.Fatalf("expected a minted key, got key=%v raw=%q", key, raw)
	}
	if key.SubjectID != caller {
		t.Fatalf("self-bound key subject = %v, want caller %v", key.SubjectID, caller)
	}
	if key.CreatedBy != caller {
		t.Fatalf("created_by = %v, want caller %v", key.CreatedBy, caller)
	}
}

func TestCreateAPIKeyImpersonationMemberForbidden(t *testing.T) {
	project := projectFixture()
	o := newTestOps(t, project, models.ProjectRoleMember)
	subject := uuid.New()

	_, _, err := o.ops.Create(ctxAs(uuid.New()), withSubject(baseInput(project.ID), subject))
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("expected forbidden for member impersonation mint, got %v", err)
	}
}

func TestCreateAPIKeyImpersonationAdminPasses(t *testing.T) {
	project := projectFixture()
	o := newTestOps(t, project, models.ProjectRoleAdmin)
	stubCreate(o)
	subject := uuid.New()

	key, _, err := o.ops.Create(ctxAs(uuid.New()), withSubject(baseInput(project.ID), subject))
	if err != nil {
		t.Fatalf("expected admin impersonation mint to pass, got %v", err)
	}
	if key.SubjectID != subject {
		t.Fatalf("key subject = %v, want %v", key.SubjectID, subject)
	}
}
