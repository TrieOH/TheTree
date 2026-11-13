// Package keys owns the key lifecycle for IdentityX: every scope — the
// platform, or a Project — always has active signing (Ed25519) and
// encryption (RSA-4096) keys with a real lifetime. The Token-lifecycle
// module verifies against keys but never creates them; projects,
// organizations, boot, and the periodic rotation worker all cross this
// seam instead of reaching into the repo.
//
// Ensure is idempotent and self-healing: it creates missing keys, rotates
// keys that are expired or carry no expiry (the legacy provisioning path
// stamped NULL), and sweeps retiring keys once their grace period has
// passed. Rotation is proactive — a key is replaced when it has less than
// one worker interval of life left, so there is never a window with no
// valid active key.
package keys

import (
	"context"
	"errors"
	"time"
	"uuid"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"
	"lib/database"
	"lib/telemetry"
)

// Config carries the key-policy knobs, resolved once at construction:
// KeyLifetime is the lifetime stamped on newly created keys, RefreshTTL
// is how long a retiring key keeps verifying after rotation (one
// refresh-token lifetime covers every artifact signed by it), and
// RotateInterval is the proactive-rotation lead window — the periodic
// worker runs on the same interval.
type Config struct {
	KeyLifetime    time.Duration
	RefreshTTL     time.Duration
	RotateInterval time.Duration
}

// Manager owns key provisioning and rotation for every scope.
type Manager struct {
	keys     ports.CryptoKeysRepo
	projects ports.ProjectRepo
	cfg      Config
	gen      func(models.CryptoKeyType) (*crypto.KeyPair, error)
	now      func() time.Time
	onChange func(context.Context)
}

// Option configures a Manager at construction.
type Option func(*Manager)

// WithClock overrides the Manager's clock so tests can pin rotation
// boundaries instead of sleeping.
func WithClock(now func() time.Time) Option {
	return func(m *Manager) { m.now = now }
}

// WithOnChange runs fn whenever a call creates or rotates a key, or
// EnsureAll's sweep retires one. It is best effort: the keys are already
// stored when fn runs, and fn cannot fail the call.
func WithOnChange(fn func(context.Context)) Option {
	return func(m *Manager) { m.onChange = fn }
}

// WithKeyGen overrides key generation — an internal seam for tests, which
// inject a cheap fixed pair instead of generating real Ed25519/RSA-4096
// keys sealed with the master key.
func WithKeyGen(gen func(models.CryptoKeyType) (*crypto.KeyPair, error)) Option {
	return func(m *Manager) { m.gen = gen }
}

func NewManager(keys ports.CryptoKeysRepo, projects ports.ProjectRepo, cfg Config, opts ...Option) *Manager {
	m := &Manager{
		keys:     keys,
		projects: projects,
		cfg:      cfg,
		gen:      func(t models.CryptoKeyType) (*crypto.KeyPair, error) { return crypto.GenerateKeyPair(string(t)) },
		now:      time.Now,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// keyTypes are the key types every scope is provisioned with. Encryption
// keys have no consumer today (OAuth secrets are sealed with the master
// key) but are kept symmetric so a future envelope-encryption consumer
// finds a ready lifecycle.
var keyTypes = []models.CryptoKeyType{models.SigningCryptoKeyType, models.EncryptionCryptoKeyType}

// Ensure makes the scope's keys healthy: for each key type it creates a
// fresh active key when none exists, rotates an active key that is
// expired or carries no expiry, and sweeps retiring keys whose grace
// period has elapsed. It is idempotent and safe to run on every boot and
// every worker tick.
func (m *Manager) Ensure(ctx context.Context, scope *uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "keys.Ensure")
	defer span.End()

	now := m.now()
	changed := false
	for _, t := range keyTypes {
		active, err := m.keys.GetActive(ctx, t, scope)
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return err
		}
		c, err := m.ensureType(ctx, scope, t, active, now)
		if err != nil {
			return err
		}
		changed = changed || c
	}
	err := m.sweep(ctx, scope, now)
	if err != nil {
		return err
	}
	if changed {
		m.notify(ctx)
	}
	return nil
}

// EnsureAll is Ensure for every scope: the platform (nil scope) and every
// Project. It reads the projects and every active key in two queries and
// sweeps all scopes in a third, so its cost does not grow with the number
// of projects; only a scope that is missing a key or due for rotation
// takes more round trips. Boot runs it before the router accepts traffic,
// and the periodic rotation worker runs it on each tick.
func (m *Manager) EnsureAll(ctx context.Context) error {
	ctx, span := telemetry.StartSpan(ctx, "keys.EnsureAll")
	defer span.End()

	projects, err := m.projects.ListAll(ctx)
	if err != nil {
		return err
	}
	active, err := m.keys.ListActive(ctx)
	if err != nil {
		return err
	}
	current := make(map[uuid.UUID]map[models.CryptoKeyType]*models.CryptoKey, len(projects)+1)
	for i := range active {
		k := &active[i]
		id := scopeOf(k.ProjectID)
		if current[id] == nil {
			current[id] = make(map[models.CryptoKeyType]*models.CryptoKey, len(keyTypes))
		}
		current[id][k.Type] = k
	}

	scopes := make([]*uuid.UUID, 0, len(projects)+1)
	scopes = append(scopes, nil)
	for i := range projects {
		scopes = append(scopes, &projects[i].ID)
	}

	now := m.now()
	changed := false
	for _, scope := range scopes {
		for _, t := range keyTypes {
			c, err := m.ensureType(ctx, scope, t, current[scopeOf(scope)][t], now)
			if err != nil {
				return err
			}
			changed = changed || c
		}
	}
	swept, err := m.keys.SweepAllRetiring(ctx, now.Add(-m.cfg.RefreshTTL))
	if err != nil {
		return err
	}
	if changed || swept > 0 {
		m.notify(ctx)
	}
	return nil
}

// ensureType makes one key type of one scope healthy, given its active key
// (nil when there is none). It reports whether it created a key.
func (m *Manager) ensureType(ctx context.Context, scope *uuid.UUID, t models.CryptoKeyType, active *models.CryptoKey, now time.Time) (bool, error) {
	if active == nil {
		return true, m.create(ctx, scope, t, now)
	}
	if !m.needsRotation(active, now) {
		return false, nil
	}
	// Retire conditionally (active → retiring): a concurrent worker that
	// already retired the key makes Retire a no-op, and we skip creation,
	// so overlapping runs cannot double-create.
	rotated, err := m.keys.Retire(ctx, active.ID, now)
	if err != nil || !rotated {
		return false, err
	}
	return true, m.create(ctx, scope, t, now)
}

// scopeOf keys a scope in a map: uuid.Nil() stands for the platform.
func scopeOf(scope *uuid.UUID) uuid.UUID {
	if scope == nil {
		return uuid.Nil()
	}
	return *scope
}

// notify tells the OnChange hook that keys were created, rotated or
// retired, so copies of the public keys held elsewhere can be dropped.
func (m *Manager) notify(ctx context.Context) {
	if m.onChange != nil {
		m.onChange(ctx)
	}
}

// needsRotation reports whether the active key must be rotated: it is
// already expired, carries no expiry (legacy provisioning — rotate it
// once so it gets a real lifetime), or has less than one RotateInterval
// of life left — the proactive lead window.
func (m *Manager) needsRotation(key *models.CryptoKey, now time.Time) bool {
	if key.ExpiresAt == nil {
		return true
	}
	return !now.Add(m.cfg.RotateInterval).Before(*key.ExpiresAt)
}

// create generates and persists a fresh active key of the type for the
// scope, stamped with the configured lifetime.
func (m *Manager) create(ctx context.Context, scope *uuid.UUID, t models.CryptoKeyType, now time.Time) error {
	pair, err := m.gen(t)
	if err != nil {
		return err
	}
	expires := now.Add(m.cfg.KeyLifetime)
	_, err = m.keys.Create(ctx, scope, pair, string(t), &expires)
	return err
}

// sweep retires the scope's keys whose grace period has elapsed: one
// RefreshTTL after rotation, no valid token signed by them remains.
func (m *Manager) sweep(ctx context.Context, scope *uuid.UUID, now time.Time) error {
	return m.keys.SweepRetiring(ctx, scope, now.Add(-m.cfg.RefreshTTL))
}
