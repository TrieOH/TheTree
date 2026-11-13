package keys

import (
	"context"
	"errors"
	"testing"
	"time"
	"uuid"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"
	"lib/database"
)

const (
	testKeyLifetime    = 168 * time.Hour
	testRefreshTTL     = 7 * time.Hour
	testRotateInterval = 1 * time.Hour
)

func testConfig() Config {
	return Config{KeyLifetime: testKeyLifetime, RefreshTTL: testRefreshTTL, RotateInterval: testRotateInterval}
}

// ── fakes ────────────────────────────────────────────────────────────────

type fakeKey struct {
	id        uuid.UUID
	scope     *uuid.UUID
	typ       models.CryptoKeyType
	status    models.CryptoKeyStatus
	expiresAt *time.Time
	rotatedAt *time.Time
}

type fakeKeys struct {
	keys []fakeKey
	// retireGives overrides Retire's faithful answer — used to simulate a
	// concurrent worker that already retired the key.
	retireGives *bool
	// perScopeReads counts the per-scope calls (GetActive, SweepRetiring),
	// which EnsureAll must not make.
	perScopeReads int
}

func sameScope(a, b *uuid.UUID) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func (f *fakeKeys) find(id uuid.UUID) *fakeKey {
	for i := range f.keys {
		if f.keys[i].id == id {
			return &f.keys[i]
		}
	}
	return nil
}

func (f *fakeKeys) toModel(k *fakeKey) *models.CryptoKey {
	return &models.CryptoKey{
		ID:                  k.id,
		ProjectID:           k.scope,
		Type:                k.typ,
		Status:              k.status,
		PublicKey:           "pub",
		EncryptedPrivateKey: "enc",
		Algorithm:           "test",
		ExpiresAt:           k.expiresAt,
		RotatedAt:           k.rotatedAt,
	}
}

func (f *fakeKeys) GetActive(_ context.Context, typ models.CryptoKeyType, scope *uuid.UUID) (*models.CryptoKey, error) {
	f.perScopeReads++
	for i := range f.keys {
		k := &f.keys[i]
		if k.typ == typ && k.status == models.CryptoKeyStatusActive && sameScope(k.scope, scope) {
			return f.toModel(k), nil
		}
	}
	return nil, database.ErrNotFound
}

func (f *fakeKeys) GetByID(_ context.Context, id uuid.UUID) (*models.CryptoKey, error) {
	if k := f.find(id); k != nil {
		return f.toModel(k), nil
	}
	return nil, database.ErrNotFound
}

func (f *fakeKeys) GetActiveSigningKeys(_ context.Context, _ *uuid.UUID) ([]models.ActiveSigningKey, error) {
	return nil, nil
}

func (f *fakeKeys) Create(_ context.Context, scope *uuid.UUID, _ *crypto.KeyPair, typ string, expiresAt *time.Time) (*models.CryptoKey, error) {
	k := fakeKey{
		id: uuid.New(), scope: scope, typ: models.CryptoKeyType(typ),
		status: models.CryptoKeyStatusActive, expiresAt: expiresAt,
	}
	f.keys = append(f.keys, k)
	return f.toModel(&f.keys[len(f.keys)-1]), nil
}

func (f *fakeKeys) Retire(_ context.Context, id uuid.UUID, rotatedAt time.Time) (bool, error) {
	if f.retireGives != nil {
		return *f.retireGives, nil
	}
	k := f.find(id)
	if k == nil || k.status != models.CryptoKeyStatusActive {
		return false, nil
	}
	k.status = models.CryptoKeyStatusRetiring
	k.rotatedAt = &rotatedAt
	return true, nil
}

func (f *fakeKeys) SweepRetiring(_ context.Context, scope *uuid.UUID, before time.Time) error {
	f.perScopeReads++
	for i := range f.keys {
		k := &f.keys[i]
		if k.status != models.CryptoKeyStatusRetiring || k.rotatedAt == nil || !sameScope(k.scope, scope) {
			continue
		}
		if k.rotatedAt.Before(before) {
			k.status = models.CryptoKeyStatusRetired
		}
	}
	return nil
}

func (f *fakeKeys) ListActive(_ context.Context) ([]models.CryptoKey, error) {
	var out []models.CryptoKey
	for i := range f.keys {
		if f.keys[i].status == models.CryptoKeyStatusActive {
			out = append(out, *f.toModel(&f.keys[i]))
		}
	}
	return out, nil
}

func (f *fakeKeys) SweepAllRetiring(_ context.Context, before time.Time) (int64, error) {
	var n int64
	for i := range f.keys {
		k := &f.keys[i]
		if k.status == models.CryptoKeyStatusRetiring && k.rotatedAt != nil && k.rotatedAt.Before(before) {
			k.status = models.CryptoKeyStatusRetired
			n++
		}
	}
	return n, nil
}

type fakeProjects struct {
	projects []models.Project
}

func (f *fakeProjects) ListAll(_ context.Context) ([]models.Project, error) { return f.projects, nil }
func (f *fakeProjects) Create(_ context.Context, _ models.Project) (*models.Project, error) {
	return nil, errors.New("stub")
}
func (f *fakeProjects) GetByID(_ context.Context, _ uuid.UUID) (*models.Project, error) {
	return nil, errors.New("stub")
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
func (f *fakeProjects) GetRole(_ context.Context, _, _ uuid.UUID) (models.ProjectRole, error) {
	return "", errors.New("stub")
}
func (f *fakeProjects) ListMembers(_ context.Context, _ uuid.UUID) ([]models.ProjectMember, error) {
	return nil, errors.New("stub")
}

var _ ports.CryptoKeysRepo = (*fakeKeys)(nil)
var _ ports.ProjectRepo = (*fakeProjects)(nil)

// ── fixture ──────────────────────────────────────────────────────────────

type fixture struct {
	mgr  *Manager
	keys *fakeKeys
	now  time.Time
	ctx  context.Context
	// changes counts the OnChange notifications.
	changes int
}

func newFixture(t *testing.T, projects []models.Project) *fixture {
	t.Helper()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	f := &fixture{
		keys: &fakeKeys{},
		now:  now,
		ctx:  context.Background(),
	}
	f.mgr = NewManager(f.keys, &fakeProjects{projects: projects}, testConfig(),
		WithClock(func() time.Time { return now }),
		WithKeyGen(func(models.CryptoKeyType) (*crypto.KeyPair, error) {
			return &crypto.KeyPair{Public: "pub", EncryptedPrivate: "enc", Algorithm: "test"}, nil
		}),
		WithOnChange(func(context.Context) { f.changes++ }),
	)
	return f
}

func (f *fixture) mustEnsure(t *testing.T, scope *uuid.UUID) {
	t.Helper()
	err := f.mgr.Ensure(f.ctx, scope)
	if err != nil {
		t.Fatalf("Ensure: %v", err)
	}
}

func (f *fixture) mustEnsureAll(t *testing.T) {
	t.Helper()
	err := f.mgr.EnsureAll(f.ctx)
	if err != nil {
		t.Fatalf("EnsureAll: %v", err)
	}
}

func (f *fixture) seed(t *testing.T, scope *uuid.UUID, expiresAt *time.Time) {
	t.Helper()
	_, err := f.keys.Create(f.ctx, scope, nil, "signing", expiresAt)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
}

func (f *fixture) active(scope *uuid.UUID, typ models.CryptoKeyType) *fakeKey {
	for i := range f.keys.keys {
		k := &f.keys.keys[i]
		if k.typ == typ && k.status == models.CryptoKeyStatusActive && sameScope(k.scope, scope) {
			return k
		}
	}
	return nil
}

func (f *fixture) count(scope *uuid.UUID, typ models.CryptoKeyType) int {
	n := 0
	for _, k := range f.keys.keys {
		if k.typ == typ && sameScope(k.scope, scope) {
			n++
		}
	}
	return n
}

func (f *fixture) expires(scope *uuid.UUID, typ models.CryptoKeyType) time.Time {
	k := f.active(scope, typ)
	if k == nil || k.expiresAt == nil {
		return time.Time{}
	}
	return *k.expiresAt
}

// ── Ensure: creation ─────────────────────────────────────────────────────

func TestEnsureCreatesBothTypesWhenMissing(t *testing.T) {
	f := newFixture(t, nil)
	scope := new(uuid.New())

	f.mustEnsure(t, scope)

	for _, typ := range []models.CryptoKeyType{models.SigningCryptoKeyType, models.EncryptionCryptoKeyType} {
		k := f.active(scope, typ)
		if k == nil {
			t.Fatalf("no active %s key after Ensure", typ)
		}
		if k.expiresAt == nil || !k.expiresAt.Equal(f.now.Add(testKeyLifetime)) {
			t.Fatalf("%s key expires %v, want %v", typ, k.expiresAt, f.now.Add(testKeyLifetime))
		}
	}
}

func TestEnsureCreatesPlatformKeysWhenMissing(t *testing.T) {
	f := newFixture(t, nil)

	f.mustEnsure(t, nil)

	if f.active(nil, models.SigningCryptoKeyType) == nil {
		t.Fatal("no platform signing key after Ensure")
	}
}

func TestEnsureIsIdempotent(t *testing.T) {
	f := newFixture(t, nil)
	scope := new(uuid.New())

	f.mustEnsure(t, scope)
	f.mustEnsure(t, scope)

	if n := f.count(scope, models.SigningCryptoKeyType); n != 1 {
		t.Fatalf("signing key count = %d after two Ensures, want 1", n)
	}
	if n := f.count(scope, models.EncryptionCryptoKeyType); n != 1 {
		t.Fatalf("encryption key count = %d after two Ensures, want 1", n)
	}
}

// ── Ensure: rotation ─────────────────────────────────────────────────────

func TestEnsureRotatesExpiredKey(t *testing.T) {
	f := newFixture(t, nil)
	scope := new(uuid.New())
	expired := f.now.Add(-time.Hour)
	f.seed(t, scope, &expired)
	old := f.active(scope, models.SigningCryptoKeyType)

	f.mustEnsure(t, scope)

	if old.status != models.CryptoKeyStatusRetiring || old.rotatedAt == nil || !old.rotatedAt.Equal(f.now) {
		t.Fatalf("old key status=%s rotated_at=%v, want retiring at %v", old.status, old.rotatedAt, f.now)
	}
	cur := f.active(scope, models.SigningCryptoKeyType)
	if cur == nil || cur.id == old.id {
		t.Fatal("no fresh active key after rotation")
	}
	if !f.expires(scope, models.SigningCryptoKeyType).Equal(f.now.Add(testKeyLifetime)) {
		t.Fatalf("fresh key expires %v, want %v", f.expires(scope, models.SigningCryptoKeyType), f.now.Add(testKeyLifetime))
	}
}

func TestEnsureRotatesLegacyNilExpiryOnce(t *testing.T) {
	f := newFixture(t, nil)
	scope := new(uuid.New())
	f.seed(t, scope, nil)

	f.mustEnsure(t, scope)

	if n := f.count(scope, models.SigningCryptoKeyType); n != 2 {
		t.Fatalf("key count = %d after first Ensure, want 2 (retiring + active)", n)
	}

	// second Ensure: the new key has a real lifetime, so nothing rotates again
	f.mustEnsure(t, scope)

	if n := f.count(scope, models.SigningCryptoKeyType); n != 2 {
		t.Fatalf("key count = %d after second Ensure, want 2 (rotation once)", n)
	}
}

func TestEnsureDoesNotRotateFreshKey(t *testing.T) {
	f := newFixture(t, nil)
	scope := new(uuid.New())
	far := f.now.Add(testKeyLifetime)
	f.seed(t, scope, &far)
	before := f.active(scope, models.SigningCryptoKeyType)

	f.mustEnsure(t, scope)

	if got := f.active(scope, models.SigningCryptoKeyType); got == nil || got.id != before.id {
		t.Fatal("fresh key was rotated when it had a full lifetime left")
	}
}

func TestEnsureRotatesWhenRemainingWithinInterval(t *testing.T) {
	f := newFixture(t, nil)
	scope := new(uuid.New())
	soon := f.now.Add(testRotateInterval - time.Minute)
	f.seed(t, scope, &soon)
	before := f.active(scope, models.SigningCryptoKeyType)

	f.mustEnsure(t, scope)

	if got := f.active(scope, models.SigningCryptoKeyType); got == nil || got.id == before.id {
		t.Fatal("key inside the lead window was not rotated")
	}
}

// ── Ensure: sweep ────────────────────────────────────────────────────────

func TestEnsureSweepsRetiringAfterGrace(t *testing.T) {
	f := newFixture(t, nil)
	scope := new(uuid.New())
	oldRotated := f.now.Add(-testRefreshTTL - time.Minute)
	f.keys.keys = append(f.keys.keys, fakeKey{
		id: uuid.New(), scope: scope, typ: models.SigningCryptoKeyType,
		status: models.CryptoKeyStatusRetiring, rotatedAt: &oldRotated,
	})
	freshRotated := f.now.Add(-testRefreshTTL + time.Minute)
	f.keys.keys = append(f.keys.keys, fakeKey{
		id: uuid.New(), scope: scope, typ: models.EncryptionCryptoKeyType,
		status: models.CryptoKeyStatusRetiring, rotatedAt: &freshRotated,
	})

	f.mustEnsure(t, scope)

	oldKey := f.keys.find(f.keys.keys[0].id)
	if oldKey.status != models.CryptoKeyStatusRetired {
		t.Fatalf("old retiring key status = %s, want retired", oldKey.status)
	}
	freshKey := f.keys.find(f.keys.keys[1].id)
	if freshKey.status != models.CryptoKeyStatusRetiring {
		t.Fatalf("fresh retiring key status = %s, want still retiring (inside grace)", freshKey.status)
	}
}

// ── Ensure: concurrency ──────────────────────────────────────────────────

func TestEnsureSkipsCreationWhenRetireLostRace(t *testing.T) {
	f := newFixture(t, nil)
	scope := new(uuid.New())
	expired := f.now.Add(-time.Hour)
	f.seed(t, scope, &expired)
	// a concurrent worker already retired the active key: Retire reports
	// false, so Ensure must not create a replacement (the other worker did)
	no := false
	f.keys.retireGives = &no

	f.mustEnsure(t, scope)

	if n := f.count(scope, models.SigningCryptoKeyType); n != 1 {
		t.Fatalf("key count = %d after lost race, want 1 (no double-create)", n)
	}
}

// ── EnsureAll ────────────────────────────────────────────────────────────

func TestEnsureAllCoversPlatformAndEveryProject(t *testing.T) {
	pid1 := new(uuid.New())
	pid2 := new(uuid.New())
	f := newFixture(t, []models.Project{{ID: *pid1}, {ID: *pid2}})

	f.mustEnsureAll(t)

	if f.active(nil, models.SigningCryptoKeyType) == nil {
		t.Fatal("platform scope has no signing key")
	}
	for _, pid := range []uuid.UUID{*pid1, *pid2} {
		if f.active(&pid, models.SigningCryptoKeyType) == nil {
			t.Fatalf("project %s has no signing key", pid)
		}
		if f.active(&pid, models.EncryptionCryptoKeyType) == nil {
			t.Fatalf("project %s has no encryption key", pid)
		}
	}
}

func TestEnsureAllMakesNoPerScopeCalls(t *testing.T) {
	projects := make([]models.Project, 0, 20)
	for range 20 {
		projects = append(projects, models.Project{ID: uuid.New()})
	}
	f := newFixture(t, projects)
	f.mustEnsureAll(t)
	f.keys.perScopeReads = 0

	f.mustEnsureAll(t)

	if f.keys.perScopeReads != 0 {
		t.Fatalf("EnsureAll made %d per-scope calls, want 0", f.keys.perScopeReads)
	}
}

func TestEnsureAllSweepsEveryScope(t *testing.T) {
	pid := new(uuid.New())
	f := newFixture(t, []models.Project{{ID: *pid}})
	f.mustEnsureAll(t)

	// A retiring key in each scope, rotated before the grace period began.
	rotated := f.now.Add(-testRefreshTTL - time.Minute)
	for _, scope := range []*uuid.UUID{nil, pid} {
		f.keys.keys = append(f.keys.keys, fakeKey{
			id: uuid.New(), scope: scope, typ: models.SigningCryptoKeyType,
			status: models.CryptoKeyStatusRetiring, rotatedAt: &rotated,
		})
	}

	f.mustEnsureAll(t)

	for _, k := range f.keys.keys {
		if k.status == models.CryptoKeyStatusRetiring {
			t.Fatalf("key %s still retiring after EnsureAll", k.id)
		}
	}
}

func TestEnsureAllNotifiesOnlyWhenKeysChange(t *testing.T) {
	f := newFixture(t, []models.Project{{ID: uuid.New()}})

	f.mustEnsureAll(t)
	if f.changes != 1 {
		t.Fatalf("changes = %d after provisioning, want 1", f.changes)
	}

	f.mustEnsureAll(t)
	if f.changes != 1 {
		t.Fatalf("changes = %d after a pass with nothing to do, want 1", f.changes)
	}

	// A sweep alone changes the published keys too.
	rotated := f.now.Add(-testRefreshTTL - time.Minute)
	f.keys.keys = append(f.keys.keys, fakeKey{
		id: uuid.New(), typ: models.SigningCryptoKeyType,
		status: models.CryptoKeyStatusRetiring, rotatedAt: &rotated,
	})
	f.mustEnsureAll(t)
	if f.changes != 2 {
		t.Fatalf("changes = %d after a sweep, want 2", f.changes)
	}
}

func TestEnsureNotifiesWhenItCreates(t *testing.T) {
	f := newFixture(t, nil)

	f.mustEnsure(t, new(uuid.New()))

	if f.changes != 1 {
		t.Fatalf("changes = %d, want 1", f.changes)
	}
}
