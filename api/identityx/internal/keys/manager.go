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
	"time"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
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
}

// Option configures a Manager at construction.
type Option func(*Manager)

// WithClock overrides the Manager's clock so tests can pin rotation
// boundaries instead of sleeping.
func WithClock(now func() time.Time) Option {
	return func(m *Manager) { m.now = now }
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
	for _, t := range keyTypes {
		active, err := m.keys.GetActive(ctx, t, scope)
		if err != nil {
			if !fun.Is(err, fun.CodeNotFound) {
				return err
			}
			err = m.create(ctx, scope, t, now)
			if err != nil {
				return err
			}
			continue
		}
		if !m.needsRotation(active, now) {
			continue
		}
		// Retire conditionally (active → retiring): a concurrent worker
		// that already retired the key makes Retire a no-op, and we skip
		// creation, so overlapping runs cannot double-create.
		rotated, err := m.keys.Retire(ctx, active.ID, now)
		if err != nil {
			return err
		}
		if rotated {
			err = m.create(ctx, scope, t, now)
			if err != nil {
				return err
			}
		}
	}
	return m.sweep(ctx, scope, now)
}

// EnsureAll makes every scope healthy: the platform (nil scope) and every
// Project. Boot calls it once (replacing the old per-project enqueue);
// the periodic rotation worker calls it on each tick.
func (m *Manager) EnsureAll(ctx context.Context) error {
	ctx, span := telemetry.StartSpan(ctx, "keys.EnsureAll")
	defer span.End()

	err := m.Ensure(ctx, nil)
	if err != nil {
		return err
	}
	projects, err := m.projects.ListAll(ctx)
	if err != nil {
		return err
	}
	for i := range projects {
		pid := projects[i].ID
		err = m.Ensure(ctx, &pid)
		if err != nil {
			return err
		}
	}
	return nil
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
