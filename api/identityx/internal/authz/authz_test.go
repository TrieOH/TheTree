package authz

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"IdentityX/models"
	"IdentityX/ports"
	libauthz "lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// ── In-memory fakes (the seam's second adapter) ──────────────────────────

type fakeProjects struct {
	projects map[uuid.UUID]*models.Project
	roles    map[[2]uuid.UUID]models.ProjectRole // {actorID, projectID}
}

func newFakeProjects() *fakeProjects {
	return &fakeProjects{projects: map[uuid.UUID]*models.Project{}, roles: map[[2]uuid.UUID]models.ProjectRole{}}
}

func (f *fakeProjects) add(project *models.Project) *fakeProjects {
	f.projects[project.ID] = project
	return f
}

func (f *fakeProjects) role(actorID, projectID uuid.UUID, r models.ProjectRole) *fakeProjects {
	f.roles[[2]uuid.UUID{actorID, projectID}] = r
	return f
}

func (f *fakeProjects) Create(_ context.Context, _ models.Project) (*models.Project, error) {
	return nil, errors.New("stub")
}
func (f *fakeProjects) GetByID(_ context.Context, id uuid.UUID) (*models.Project, error) {
	p, ok := f.projects[id]
	if !ok {
		return nil, fun.Err("project not found").NotFound()
	}
	return p, nil
}
func (f *fakeProjects) ListByOrganization(_ context.Context, _ uuid.UUID) ([]models.Project, error) {
	return nil, errors.New("stub")
}
func (f *fakeProjects) ListJoined(_ context.Context, _ uuid.UUID) ([]models.Project, error) {
	return nil, errors.New("stub")
}
func (f *fakeProjects) ListOwned(_ context.Context, _ uuid.UUID) ([]models.Project, error) {
	return nil, errors.New("stub")
}
func (f *fakeProjects) AddMember(_ context.Context, _ models.ProjectMember) error {
	return errors.New("stub")
}
func (f *fakeProjects) RemoveMember(_ context.Context, _, _ uuid.UUID) error {
	return errors.New("stub")
}
func (f *fakeProjects) GetMember(_ context.Context, _, _ uuid.UUID) (*models.ProjectMember, error) {
	return nil, errors.New("stub")
}
func (f *fakeProjects) GetRole(_ context.Context, actorID, projectID uuid.UUID) (models.ProjectRole, error) {
	r, ok := f.roles[[2]uuid.UUID{actorID, projectID}]
	if !ok {
		return "", fun.Err("project member not found").NotFound()
	}
	return r, nil
}
func (f *fakeProjects) ListMembers(_ context.Context, _ uuid.UUID) ([]models.ProjectMember, error) {
	return nil, errors.New("stub")
}

type fakeOrgs struct {
	roles map[[2]uuid.UUID]models.OrganizationRole // {actorID, orgID}
}

func newFakeOrgs() *fakeOrgs { return &fakeOrgs{roles: map[[2]uuid.UUID]models.OrganizationRole{}} }

func (f *fakeOrgs) role(actorID, orgID uuid.UUID, r models.OrganizationRole) *fakeOrgs {
	f.roles[[2]uuid.UUID{actorID, orgID}] = r
	return f
}

func (f *fakeOrgs) Create(_ context.Context, _ models.Organization) (*models.Organization, error) {
	return nil, errors.New("stub")
}
func (f *fakeOrgs) GetByID(_ context.Context, _ uuid.UUID) (*models.Organization, error) {
	return nil, errors.New("stub")
}
func (f *fakeOrgs) ListOwned(_ context.Context, _ uuid.UUID) ([]models.Organization, error) {
	return nil, errors.New("stub")
}
func (f *fakeOrgs) ListJoined(_ context.Context, _ uuid.UUID) ([]models.Organization, error) {
	return nil, errors.New("stub")
}
func (f *fakeOrgs) AddMember(_ context.Context, _ models.OrganizationMember) error {
	return errors.New("stub")
}
func (f *fakeOrgs) RemoveMember(_ context.Context, _, _ uuid.UUID) error {
	return errors.New("stub")
}
func (f *fakeOrgs) GetMember(_ context.Context, _, _ uuid.UUID) (*models.OrganizationMember, error) {
	return nil, errors.New("stub")
}
func (f *fakeOrgs) GetRole(_ context.Context, actorID, orgID uuid.UUID) (models.OrganizationRole, error) {
	r, ok := f.roles[[2]uuid.UUID{actorID, orgID}]
	if !ok {
		return "", fun.Err("org member not found").NotFound()
	}
	return r, nil
}
func (f *fakeOrgs) ListMembers(_ context.Context, _ uuid.UUID) ([]models.OrganizationMember, error) {
	return nil, errors.New("stub")
}

type fakePlatformRoles struct {
	roles map[uuid.UUID]models.PlatformRole
}

func newFakePlatformRoles() *fakePlatformRoles {
	return &fakePlatformRoles{roles: map[uuid.UUID]models.PlatformRole{}}
}

func (f *fakePlatformRoles) role(actorID uuid.UUID, r models.PlatformRole) *fakePlatformRoles {
	f.roles[actorID] = r
	return f
}

func (f *fakePlatformRoles) Give(_ context.Context, _ uuid.UUID, _ models.PlatformRole, _ *json.RawMessage) (*models.PlatformRoleRelation, error) {
	return nil, errors.New("stub")
}
func (f *fakePlatformRoles) GetRole(_ context.Context, actorID uuid.UUID) (models.PlatformRole, error) {
	r, ok := f.roles[actorID]
	if !ok {
		return "", fun.Err("platform role not found").NotFound()
	}
	return r, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────

var (
	user1 = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	user2 = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	org1  = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	proj1 = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	proj2 = uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
)

func newService(projects ports.ProjectRepo, orgs ports.OrganizationRepo, platform ports.PlatformRolesRepo) *Service {
	return New(orgs, projects, platform)
}

func forbidden(err error) bool {
	return err != nil && fun.Is(err, fun.CodeForbidden)
}

// ── CheckProject: direct project membership ──────────────────────────────

func TestCheckProjectDirectMemberPasses(t *testing.T) {
	s := newService(newFakeProjects().add(&models.Project{ID: proj1}).role(user1, proj1, models.ProjectRoleAdmin), newFakeOrgs(), newFakePlatformRoles())
	err := s.CheckProject(context.Background(), user1, proj1, models.ProjectRoleAdmin)
	if err != nil {
		t.Fatalf("expected admin to pass, got %v", err)
	}
}

func TestCheckProjectDirectMemberInsufficientRole(t *testing.T) {
	s := newService(newFakeProjects().add(&models.Project{ID: proj1}).role(user1, proj1, models.ProjectRoleMember), newFakeOrgs(), newFakePlatformRoles())
	err := s.CheckProject(context.Background(), user1, proj1, models.ProjectRoleAdmin)
	if !forbidden(err) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

// ── CheckProject: org fallback ───────────────────────────────────────────

func TestCheckProjectOrgMemberFallsBack(t *testing.T) {
	// user1 has no project_members row but is an org member — the org role
	// casts onto the project (member → member).
	s := newService(
		newFakeProjects().add(&models.Project{ID: proj1, OrganizationID: &org1}),
		newFakeOrgs().role(user1, org1, models.OrganizationRoleMember),
		newFakePlatformRoles(),
	)
	err := s.CheckProject(context.Background(), user1, proj1, models.ProjectRoleMember)
	if err != nil {
		t.Fatalf("expected org member to fall back, got %v", err)
	}
}

func TestCheckProjectOrgOwnerFallsBackToOwner(t *testing.T) {
	// the org owner who created the project holds no project_members row —
	// the original lockout bug; the org owner must be project-owner-equivalent.
	s := newService(
		newFakeProjects().add(&models.Project{ID: proj1, OrganizationID: &org1}),
		newFakeOrgs().role(user1, org1, models.OrganizationRoleOwner),
		newFakePlatformRoles(),
	)
	err := s.CheckProject(context.Background(), user1, proj1, models.ProjectRoleAdmin)
	if err != nil {
		t.Fatalf("expected org owner to pass admin check via fallback, got %v", err)
	}
}

func TestCheckProjectOrgMemberInsufficientCastRole(t *testing.T) {
	// org member cast = project member (rank 0) < admin (rank 1) → forbidden.
	s := newService(
		newFakeProjects().add(&models.Project{ID: proj1, OrganizationID: &org1}),
		newFakeOrgs().role(user1, org1, models.OrganizationRoleMember),
		newFakePlatformRoles(),
	)
	err := s.CheckProject(context.Background(), user1, proj1, models.ProjectRoleAdmin)
	if !forbidden(err) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestCheckProjectDirectRoleBeatsOrgRole(t *testing.T) {
	// a project_members row wins over the org fallback — both paths grant
	// member access here, and neither must conflict.
	s := newService(
		newFakeProjects().add(&models.Project{ID: proj1, OrganizationID: &org1}).role(user1, proj1, models.ProjectRoleMember),
		newFakeOrgs().role(user1, org1, models.OrganizationRoleMember),
		newFakePlatformRoles(),
	)
	err := s.CheckProject(context.Background(), user1, proj1, models.ProjectRoleMember)
	if err != nil {
		t.Fatalf("expected pass, got %v", err)
	}
}

// ── CheckProject: non-members and unknown projects ───────────────────────

func TestCheckProjectNonMemberInOrgProjectForbidden(t *testing.T) {
	s := newService(
		newFakeProjects().add(&models.Project{ID: proj1, OrganizationID: &org1}),
		newFakeOrgs(), // user2 has no org role
		newFakePlatformRoles(),
	)
	err := s.CheckProject(context.Background(), user2, proj1, models.ProjectRoleMember)
	if !forbidden(err) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestCheckProjectNonMemberOrglessProjectForbidden(t *testing.T) {
	// a platform-created project has no org — no fallback, pure member check.
	s := newService(
		newFakeProjects().add(&models.Project{ID: proj2}), // OrganizationID nil
		newFakeOrgs(),
		newFakePlatformRoles(),
	)
	err := s.CheckProject(context.Background(), user2, proj2, models.ProjectRoleMember)
	if !forbidden(err) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestCheckProjectUnknownProjectNotFound(t *testing.T) {
	s := newService(newFakeProjects(), newFakeOrgs(), newFakePlatformRoles())
	err := s.CheckProject(context.Background(), user1, proj1, models.ProjectRoleMember)
	if err == nil || !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("expected NotFound, got %v", err)
	}
}

// ── CheckOrg ─────────────────────────────────────────────────────────────

func TestCheckOrgPassesAndRejects(t *testing.T) {
	s := newService(newFakeProjects(), newFakeOrgs().role(user1, org1, models.OrganizationRoleAdmin), newFakePlatformRoles())
	err := s.CheckOrg(context.Background(), user1, org1, models.OrganizationRoleAdmin)
	if err != nil {
		t.Fatalf("expected admin to pass, got %v", err)
	}
	err = s.CheckOrg(context.Background(), user2, org1, models.OrganizationRoleMember)
	if !forbidden(err) {
		t.Fatalf("expected forbidden for non-member, got %v", err)
	}
}

// ── CheckPlatform ────────────────────────────────────────────────────────

func TestCheckPlatformTiers(t *testing.T) {
	s := newService(
		newFakeProjects(),
		newFakeOrgs(),
		newFakePlatformRoles().
			role(user1, models.PlatformRoleSuperAdmin).
			role(user2, models.PlatformRoleSupport),
	)
	err := s.CheckPlatform(context.Background(), user1, models.PlatformRoleAdmin)
	if err != nil {
		t.Fatalf("expected super_admin to pass admin check, got %v", err)
	}
	err = s.CheckPlatform(context.Background(), user2, models.PlatformRoleAdmin)
	if !forbidden(err) {
		t.Fatalf("expected support to be forbidden for admin check, got %v", err)
	}
	err = s.CheckPlatform(context.Background(), uuid.New(), models.PlatformRoleSupport)
	if !forbidden(err) {
		t.Fatalf("expected unknown actor to be forbidden, got %v", err)
	}
}

var _ libauthz.Role = models.PlatformRoleAdmin // PlatformRole satisfies the shared Role interface
