package tos

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/internal/emails"
	"IdentityX/internal/repos"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"lib/database"
	"lib/testdb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// fakeEnqueuer captures the notification fan-out without a River client.
type fakeEnqueuer struct {
	args []emails.SendTosUpdateArgs
}

func (f *fakeEnqueuer) Insert(_ context.Context, args river.JobArgs, _ *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	f.args = append(f.args, args.(emails.SendTosUpdateArgs))
	return &rivertype.JobInsertResult{}, nil
}

var _ emails.Enqueuer = (*fakeEnqueuer)(nil)

// tosFixture is a fully wired tos Operations over a real Postgres, with an
// owner (project admin) and a plain project user.
type tosFixture struct {
	ops       *Operations
	pool      *pgxpool.Pool
	repos     *repos.Repos
	enqueuer  *fakeEnqueuer
	projectID uuid.UUID
	userID    uuid.UUID
	userEmail string
	ownerCtx  context.Context
	userCtx   context.Context
}

func setupTosFixture(t *testing.T) *tosFixture {
	t.Helper()

	pool := testdb.Postgres(t, "../../../db/migrations")
	database.SetDefaultRunner(database.NewPGXTxRunner(pool))

	ctx := context.Background()
	r := repos.New(sqlc.New(pool))
	enqueuer := &fakeEnqueuer{}

	email := "tos-owner@trieoh.com"
	owner, err := r.Actors.Register(ctx, models.Actor{
		AuthMethod:   models.PasswordAuthMethod,
		Email:        &email,
		Type:         models.HumanActorType,
		PasswordHash: new("x"),
	})
	if err != nil {
		t.Fatalf("register owner: %v", err)
	}
	domain := "https://acme.com.br"
	project, err := r.Projects.Create(ctx, models.Project{
		OwnerID:   owner.ID,
		Name:      "Acme",
		BrandSlug: "acmetos",
		Domain:    &domain,
		Metadata:  json.RawMessage("{}"),
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	// Creating via the repo skips the service-level owner membership, so
	// add it explicitly — CheckProject resolves the admin through it.
	err = r.Projects.AddMember(ctx, models.ProjectMember{
		ProjectID: project.ID,
		ActorID:   owner.ID,
		Role:      models.ProjectRoleOwner,
		Metadata:  json.RawMessage("{}"),
	})
	if err != nil {
		t.Fatalf("add owner member: %v", err)
	}
	// A plain project user (no membership row) who will accept the terms.
	userEmail := "user@acme.com.br"
	user, err := r.Actors.Register(ctx, models.Actor{
		ProjectID:    &project.ID,
		AuthMethod:   models.PasswordAuthMethod,
		Email:        &userEmail,
		Type:         models.HumanActorType,
		PasswordHash: new("x"),
	})
	if err != nil {
		t.Fatalf("register user: %v", err)
	}

	authzSvc := authz.New(r.Organizations, r.Projects, r.PlatformRoles)
	ownerCtx := models.WithIdentity(ctx, &models.Identity{
		Sub:  models.Subject{ID: owner.ID, Type: models.HumanActorType},
		Cred: models.Credential{Type: models.TokenCredentialType},
	})
	userCtx := models.WithIdentity(ctx, &models.Identity{
		Sub:  models.Subject{ID: user.ID, ProjectID: &project.ID, Type: models.HumanActorType},
		Cred: models.Credential{Type: models.TokenCredentialType},
	})

	return &tosFixture{
		ops:       NewOperations(r.Tos, r.Projects, r.Actors, authzSvc, emails.NewTosNotifier(enqueuer)),
		pool:      pool,
		repos:     r,
		enqueuer:  enqueuer,
		projectID: project.ID,
		userID:    user.ID,
		userEmail: userEmail,
		ownerCtx:  ownerCtx,
		userCtx:   userCtx,
	}
}

// TestTosLifecycleRealDB pins the terms flow end to end against a real
// Postgres (testcontainers): the v0 → v1 migration semantics, version
// bumping, the acceptance ledger's per-version uniqueness and idempotency,
// the continued-use stamp, and the notification fan-out for every created
// or updated version.
func TestTosLifecycleRealDB(t *testing.T) {
	f := setupTosFixture(t)

	// v0: no terms, GetCurrent misses.
	_, err := f.ops.Get(context.Background(), f.projectID)
	if err == nil {
		t.Fatal("project without terms must surface as not-found")
	}

	v1 := f.introduceV1(t)
	f.acceptV1(t, v1)
	f.updateToV2(t)
	f.stampContinuedUse(t)
	f.checkLedger(t)
	f.checkOutsider(t)
}

// introduceV1 covers the v0 → v1 migration: version 1 with a ~30-day
// notice window and one notification for the existing users.
func (f *tosFixture) introduceV1(t *testing.T) *models.TermsOfService {
	t.Helper()

	v1, err := f.ops.Create(f.ownerCtx, f.projectID, "Termos v1 em Português")
	if err != nil {
		t.Fatalf("create v1: %v", err)
	}
	if v1.Version != 1 {
		t.Fatalf("first version must be 1, got %d", v1.Version)
	}
	if until := time.Until(v1.EffectiveAt); until < 29*24*time.Hour || until > 31*24*time.Hour {
		t.Fatalf("effective date must sit ~30 days out, got %v", until)
	}
	if len(f.enqueuer.args) != 1 || f.enqueuer.args[0].TosVersion != 1 || f.enqueuer.args[0].ProjectID != f.projectID {
		t.Fatalf("v1 introduction must enqueue one notification, got %+v", f.enqueuer.args)
	}

	// Creating again conflicts; the admin must go through Update.
	_, err = f.ops.Create(f.ownerCtx, f.projectID, "dup")
	if err == nil {
		t.Fatal("second Create must conflict")
	}
	return v1
}

// acceptV1 covers the explicit clickwrap path: the ledger row pins the
// version with a content hash, and re-accepting is idempotent.
func (f *tosFixture) acceptV1(t *testing.T, v1 *models.TermsOfService) {
	t.Helper()

	acc1, err := f.ops.Accept(f.userCtx, f.projectID)
	if err != nil {
		t.Fatalf("accept v1: %v", err)
	}
	if acc1.TosVersion != v1.Version || len(acc1.ContentHash) != 64 {
		t.Fatalf("acceptance must pin version 1 and a sha256 hash, got %+v", acc1)
	}
	if acc1.Source != models.TosAcceptanceClickwrap {
		t.Fatalf("explicit acceptance must be clickwrap, got %q", acc1.Source)
	}
	_, err = f.ops.Accept(f.userCtx, f.projectID)
	if err != nil {
		t.Fatalf("re-accept v1 must be idempotent: %v", err)
	}
}

// updateToV2 covers the version bump: v2, a fresh notice window, and a
// second notification.
func (f *tosFixture) updateToV2(t *testing.T) *models.TermsOfService {
	t.Helper()

	time.Sleep(10 * time.Millisecond)
	v2, err := f.ops.Update(f.ownerCtx, f.projectID, "Termos v2 em Português")
	if err != nil {
		t.Fatalf("update to v2: %v", err)
	}
	if v2.Version != 2 {
		t.Fatalf("update must bump to version 2, got %d", v2.Version)
	}
	if len(f.enqueuer.args) != 2 || f.enqueuer.args[1].TosVersion != 2 {
		t.Fatalf("v2 update must enqueue one notification, got %+v", f.enqueuer.args)
	}
	return v2
}

// stampContinuedUse covers the login stamp: with v2 already effective, the
// first StampUse records a continued_use row and the next one is a no-op.
func (f *tosFixture) stampContinuedUse(t *testing.T) {
	t.Helper()

	// The test pulls the effective date into the past to skip the notice
	// window (the login path delegates to StampUse).
	_, err := f.pool.Exec(context.Background(), "UPDATE terms_of_service SET effective_at = NOW() - INTERVAL '1 hour' WHERE project_id = $1", f.projectID)
	if err != nil {
		t.Fatalf("pull effective_at into the past: %v", err)
	}
	err = f.ops.StampUse(context.Background(), f.userID, f.projectID)
	if err != nil {
		t.Fatalf("stamp continued use: %v", err)
	}
	// Stamping again (next login) must not duplicate the row.
	err = f.ops.StampUse(context.Background(), f.userID, f.projectID)
	if err != nil {
		t.Fatalf("re-stamp must be idempotent: %v", err)
	}
}

// checkLedger verifies the admin audit listing shows exactly one row per
// accepted version with the right consent source and the actor email.
func (f *tosFixture) checkLedger(t *testing.T) {
	t.Helper()

	listing, err := f.ops.ListAcceptances(f.ownerCtx, f.projectID)
	if err != nil {
		t.Fatalf("list acceptances: %v", err)
	}
	if len(listing) != 2 {
		t.Fatalf("ledger must hold one row per accepted version, got %d", len(listing))
	}
	sawClickwrap, sawContinuedUse := false, false
	for _, row := range listing {
		if row.ActorEmail == nil || *row.ActorEmail != f.userEmail {
			continue
		}
		switch {
		case row.TosVersion == 1 && row.Source == models.TosAcceptanceClickwrap:
			sawClickwrap = true
		case row.TosVersion == 2 && row.Source == models.TosAcceptanceContinuedUse:
			sawContinuedUse = true
		}
	}
	if !sawClickwrap || !sawContinuedUse {
		t.Fatalf("ledger must show clickwrap v1 and continued_use v2, got %+v", listing)
	}
}

// checkOutsider verifies a platform-level actor outside the project can
// never record an acceptance.
func (f *tosFixture) checkOutsider(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	outsiderEmail := "outsider@trieoh.com"
	outsider, err := f.repos.Actors.Register(ctx, models.Actor{
		AuthMethod:   models.PasswordAuthMethod,
		Email:        &outsiderEmail,
		Type:         models.HumanActorType,
		PasswordHash: new("x"),
	})
	if err != nil {
		t.Fatalf("register outsider: %v", err)
	}
	outsiderCtx := models.WithIdentity(ctx, &models.Identity{
		Sub:  models.Subject{ID: outsider.ID, Type: models.HumanActorType},
		Cred: models.Credential{Type: models.TokenCredentialType},
	})
	_, err = f.ops.Accept(outsiderCtx, f.projectID)
	if err == nil {
		t.Fatal("actor outside the project must not be able to accept")
	}
}
