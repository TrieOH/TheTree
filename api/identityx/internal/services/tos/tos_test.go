package tos

import (
	"context"
	"testing"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/internal/emails"
	"IdentityX/models"
	"IdentityX/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

// adminCtx carries an admin identity so authorizeAdmin passes.
func adminCtx() context.Context {
	actor := models.Actor{ID: uuid.New(), Type: models.HumanActorType}
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub:  models.Subject{ID: actor.ID, Type: models.HumanActorType},
		Cred: models.Credential{Type: models.TokenCredentialType},
	})
}

// userCtx carries a project-user identity (no membership row) so Accept
// resolves the caller.
func userCtx(actorID uuid.UUID) context.Context {
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub:  models.Subject{ID: actorID, Type: models.HumanActorType},
		Cred: models.Credential{Type: models.TokenCredentialType},
	})
}

func stubProject(projects ports.ProjectRepo, projectID uuid.UUID) {
	_ = mock.When(projects.GetByID(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(&models.Project{ID: projectID, Domain: new("https://acme.com.br")}, nil)
}

func stubProjectAdmin(projectID uuid.UUID) (*Operations, ports.TosRepo, ports.TosNotifier) {
	projects := mock.Mock[ports.ProjectRepo]()
	stubProject(projects, projectID)
	_ = mock.When(projects.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Equal(projectID))).
		ThenReturn(models.ProjectRoleAdmin, nil)

	tosRepo := mock.Mock[ports.TosRepo]()
	notifier := mock.Mock[ports.TosNotifier]()
	ops := NewOperations(
		tosRepo,
		projects,
		mock.Mock[ports.ActorRepo](),
		authz.New(mock.Mock[ports.OrganizationRepo](), projects, mock.Mock[ports.PlatformRolesRepo]()),
		notifier,
	)
	return ops, tosRepo, notifier
}

// stubNotify captures the ToS the notifier was called with.
func stubNotify(notifier ports.TosNotifier) func() *models.TermsOfService {
	var captured *models.TermsOfService
	_ = mock.When(notifier.NotifyUpdate(mock.AnyContext(), mock.Any[*models.Project](), mock.Any[*models.TermsOfService]())).
		ThenAnswer(func(args []any) []any {
			captured = args[2].(*models.TermsOfService)
			return []any{nil}
		})
	return func() *models.TermsOfService { return captured }
}

func TestCreateIntroducesV1AndNotifies(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	ops, tosRepo, notifier := stubProjectAdmin(projectID)

	// v0: the project has no terms yet.
	_ = mock.When(tosRepo.GetCurrent(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(nil, fun.ErrNotFound("terms_of_service not found"))
	var created *models.TermsOfService
	_ = mock.When(tosRepo.Create(mock.AnyContext(), mock.Any[models.TermsOfService]())).
		ThenAnswer(func(args []any) []any {
			c := args[1].(models.TermsOfService)
			created = &c
			return []any{created, nil}
		})
	notified := stubNotify(notifier)

	out, err := ops.Create(adminCtx(), projectID, "<h1>Termos v1</h1>")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if out.Version != 1 {
		t.Fatalf("first version must be 1, got %d", out.Version)
	}
	if until := time.Until(out.EffectiveAt); until < 29*24*time.Hour || until > 31*24*time.Hour {
		t.Fatalf("effective date must sit ~30 days out, got %v", until)
	}
	// The v0→v1 migration notifies the existing users exactly like an update.
	_ = mock.Verify(notifier, mock.Times(1)).NotifyUpdate(mock.AnyContext(), mock.Any[*models.Project](), mock.Any[*models.TermsOfService]())
	if created == nil {
		t.Fatal("repo Create must be called")
	}
	if notified() == nil || notified().Version != 1 {
		t.Fatal("the notification must carry the new version")
	}
}

func TestCreateRejectsWhenTermsAlreadyExist(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	ops, tosRepo, notifier := stubProjectAdmin(projectID)

	_ = mock.When(tosRepo.GetCurrent(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(&models.TermsOfService{ID: uuid.New(), ProjectID: projectID, Version: 1}, nil)

	_, err := ops.Create(adminCtx(), projectID, "new content")
	if err == nil || !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("creating over existing terms must conflict, got %v", err)
	}
	_ = mock.Verify(notifier, mock.Times(0)).NotifyUpdate(mock.AnyContext(), mock.Any[*models.Project](), mock.Any[*models.TermsOfService]())
}

func TestCreateRejectsEmptyContent(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	ops, tosRepo, _ := stubProjectAdmin(projectID)

	_ = mock.When(tosRepo.GetCurrent(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(nil, fun.ErrNotFound("terms_of_service not found"))

	_, err := ops.Create(adminCtx(), projectID, "")
	if err == nil || !fun.Is(err, fun.CodeValidation) {
		t.Fatalf("empty content must be rejected, got %v", err)
	}
}

func TestUpdateBumpsVersionAndNotifies(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	ops, tosRepo, notifier := stubProjectAdmin(projectID)

	var updated *models.TermsOfService
	_ = mock.When(tosRepo.Update(mock.AnyContext(), mock.Equal(projectID), mock.Any[string](), mock.Any[time.Time]())).
		ThenAnswer(func(args []any) []any {
			updated = &models.TermsOfService{
				ID: uuid.New(), ProjectID: projectID, Version: 3,
				Content: args[2].(string), EffectiveAt: args[3].(time.Time),
			}
			return []any{updated, nil}
		})
	notified := stubNotify(notifier)

	out, err := ops.Update(adminCtx(), projectID, "<h1>Termos v3</h1>")
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if out.Version != 3 {
		t.Fatalf("repo version must flow through, got %d", out.Version)
	}
	_ = mock.Verify(notifier, mock.Times(1)).NotifyUpdate(mock.AnyContext(), mock.Any[*models.Project](), mock.Any[*models.TermsOfService]())
	if notified() == nil || notified().Version != 3 {
		t.Fatal("the notification must carry the bumped version")
	}
}

func TestCurrentTosNilWhenAbsent(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	ops, tosRepo, _ := stubProjectAdmin(projectID)

	_ = mock.When(tosRepo.GetCurrent(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(nil, fun.ErrNotFound("terms_of_service not found"))

	tos, hasTos, err := ops.CurrentTos(context.Background(), projectID)
	if err != nil || hasTos || tos != nil {
		t.Fatalf("v0 project must resolve to no terms, got %v, %v, %v", tos, hasTos, err)
	}
}

func TestAcceptRecordsCurrentVersionWithContentHash(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	actorID := uuid.New()
	projects := mock.Mock[ports.ProjectRepo]()
	stubProject(projects, projectID)

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actorID))).
		ThenReturn(&models.Actor{ID: actorID, ProjectID: &projectID}, nil)

	tosRepo := mock.Mock[ports.TosRepo]()
	current := &models.TermsOfService{ID: uuid.New(), ProjectID: projectID, Version: 2, Content: "termos v2"}
	_ = mock.When(tosRepo.GetCurrent(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(current, nil)
	var recorded models.TosAcceptance
	_ = mock.When(tosRepo.RecordAcceptance(mock.AnyContext(), mock.Any[models.TosAcceptance]())).
		ThenAnswer(func(args []any) []any {
			recorded = args[1].(models.TosAcceptance)
			return []any{&recorded, nil}
		})

	ops := NewOperations(tosRepo, projects, actors,
		authz.New(mock.Mock[ports.OrganizationRepo](), projects, mock.Mock[ports.PlatformRolesRepo]()),
		emails.NewTosNotifier(mock.Mock[emails.Enqueuer]()),
	)

	out, err := ops.Accept(userCtx(actorID), projectID)
	if err != nil {
		t.Fatalf("Accept: %v", err)
	}
	if recorded.TosVersion != 2 {
		t.Fatalf("acceptance must pin the current version, got %d", recorded.TosVersion)
	}
	if recorded.Source != models.TosAcceptanceClickwrap {
		t.Fatalf("explicit acceptance must be clickwrap, got %q", recorded.Source)
	}
	if recorded.ContentHash == "" || len(recorded.ContentHash) != 64 {
		t.Fatalf("acceptance must pin a sha256 content hash, got %q", recorded.ContentHash)
	}
	if out.ActorID != actorID || out.ProjectID != projectID {
		t.Fatal("acceptance must belong to the caller and the project")
	}
}

func TestAcceptRejectsActorOutsideProject(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	actorID := uuid.New()
	otherProject := uuid.New()

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.GetByID(mock.AnyContext(), mock.Equal(actorID))).
		ThenReturn(&models.Actor{ID: actorID, ProjectID: &otherProject}, nil)

	projects := mock.Mock[ports.ProjectRepo]()
	stubProject(projects, projectID)

	ops := NewOperations(mock.Mock[ports.TosRepo](), projects, actors,
		authz.New(mock.Mock[ports.OrganizationRepo](), projects, mock.Mock[ports.PlatformRolesRepo]()),
		emails.NewTosNotifier(mock.Mock[emails.Enqueuer]()),
	)

	_, err := ops.Accept(userCtx(actorID), projectID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("actor outside the project must be rejected, got %v", err)
	}
}

func TestUpdateRequiresAdmin(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	projects := mock.Mock[ports.ProjectRepo]()
	stubProject(projects, projectID)
	_ = mock.When(projects.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Equal(projectID))).
		ThenReturn(models.ProjectRoleMember, nil)

	ops := NewOperations(mock.Mock[ports.TosRepo](), projects, mock.Mock[ports.ActorRepo](),
		authz.New(mock.Mock[ports.OrganizationRepo](), projects, mock.Mock[ports.PlatformRolesRepo]()),
		emails.NewTosNotifier(mock.Mock[emails.Enqueuer]()),
	)

	_, err := ops.Update(userCtx(uuid.New()), projectID, "content")
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("non-admin update must be forbidden, got %v", err)
	}
}

func stubStampOps(t *testing.T, projectID uuid.UUID, current *models.TermsOfService) (*Operations, ports.TosRepo) {
	t.Helper()
	projects := mock.Mock[ports.ProjectRepo]()
	stubProject(projects, projectID)
	tosRepo := mock.Mock[ports.TosRepo]()
	if current != nil {
		_ = mock.When(tosRepo.GetCurrent(mock.AnyContext(), mock.Equal(projectID))).
			ThenReturn(current, nil)
	} else {
		_ = mock.When(tosRepo.GetCurrent(mock.AnyContext(), mock.Equal(projectID))).
			ThenReturn(nil, fun.ErrNotFound("terms_of_service not found"))
	}
	return NewOperations(tosRepo, projects, mock.Mock[ports.ActorRepo](),
		authz.New(mock.Mock[ports.OrganizationRepo](), projects, mock.Mock[ports.PlatformRolesRepo]()),
		mock.Mock[ports.TosNotifier]()), tosRepo
}

func TestStampUseRecordsContinuedUseAfterEffectiveDate(t *testing.T) {
	mock.SetUp(t)
	projectID, actorID := uuid.New(), uuid.New()
	current := &models.TermsOfService{ID: uuid.New(), ProjectID: projectID, Version: 2,
		Content: "termos v2", EffectiveAt: time.Now().Add(-24 * time.Hour)}
	ops, tosRepo := stubStampOps(t, projectID, current)

	_ = mock.When(tosRepo.LatestAcceptedVersion(mock.AnyContext(), mock.Equal(actorID), mock.Equal(projectID))).
		ThenReturn(1, nil)
	var recorded models.TosAcceptance
	_ = mock.When(tosRepo.RecordAcceptance(mock.AnyContext(), mock.Any[models.TosAcceptance]())).
		ThenAnswer(func(args []any) []any {
			recorded = args[1].(models.TosAcceptance)
			return []any{&recorded, nil}
		})

	err := ops.StampUse(context.Background(), actorID, projectID)
	if err != nil {
		t.Fatalf("StampUse: %v", err)
	}
	if recorded.TosVersion != 2 {
		t.Fatalf("stamp must pin the current version, got %d", recorded.TosVersion)
	}
	if recorded.Source != models.TosAcceptanceContinuedUse {
		t.Fatalf("stamp must be continued_use, got %q", recorded.Source)
	}
	if len(recorded.ContentHash) != 64 {
		t.Fatalf("stamp must pin a sha256 content hash, got %q", recorded.ContentHash)
	}
}

func TestStampUseSkipsInsideNoticeWindow(t *testing.T) {
	mock.SetUp(t)
	projectID, actorID := uuid.New(), uuid.New()
	current := &models.TermsOfService{ID: uuid.New(), ProjectID: projectID, Version: 2,
		Content: "termos v2", EffectiveAt: time.Now().Add(29 * 24 * time.Hour)}
	ops, tosRepo := stubStampOps(t, projectID, current)

	err := ops.StampUse(context.Background(), actorID, projectID)
	if err != nil {
		t.Fatalf("StampUse must not fail inside the notice window: %v", err)
	}
	_, _ = mock.Verify(tosRepo, mock.Times(0)).RecordAcceptance(mock.AnyContext(), mock.Any[models.TosAcceptance]())
}

func TestStampUseSkipsWhenAlreadyAccepted(t *testing.T) {
	mock.SetUp(t)
	projectID, actorID := uuid.New(), uuid.New()
	current := &models.TermsOfService{ID: uuid.New(), ProjectID: projectID, Version: 2,
		Content: "termos v2", EffectiveAt: time.Now().Add(-24 * time.Hour)}
	ops, tosRepo := stubStampOps(t, projectID, current)

	_ = mock.When(tosRepo.LatestAcceptedVersion(mock.AnyContext(), mock.Equal(actorID), mock.Equal(projectID))).
		ThenReturn(2, nil)

	err := ops.StampUse(context.Background(), actorID, projectID)
	if err != nil {
		t.Fatalf("StampUse: %v", err)
	}
	_, _ = mock.Verify(tosRepo, mock.Times(0)).RecordAcceptance(mock.AnyContext(), mock.Any[models.TosAcceptance]())
}

func TestStampUseSkipsForProjectWithoutTerms(t *testing.T) {
	mock.SetUp(t)
	projectID, actorID := uuid.New(), uuid.New()
	ops, tosRepo := stubStampOps(t, projectID, nil)

	err := ops.StampUse(context.Background(), actorID, projectID)
	if err != nil {
		t.Fatalf("StampUse: %v", err)
	}
	_, _ = mock.Verify(tosRepo, mock.Times(0)).RecordAcceptance(mock.AnyContext(), mock.Any[models.TosAcceptance]())
}
