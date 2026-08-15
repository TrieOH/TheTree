package ports

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"time"

	"github.com/google/uuid"
)

type CryptoKeysRepo interface {
	GetActive(ctx context.Context, keyType models.CryptoKeyType, projectID *uuid.UUID) (*models.CryptoKey, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.CryptoKey, error)
	GetActiveSigningKeys(ctx context.Context, projectID *uuid.UUID) ([]models.ActiveSigningKey, error)
	// Create persists a fresh key for the scope. expiresAt stamps the key's
	// lifetime (nil stamps none); the Key-lifecycle module is the only
	// caller and always passes a real expiry.
	Create(ctx context.Context, projectID *uuid.UUID, pair *crypto.KeyPair, keyType string, expiresAt *time.Time) (*models.CryptoKey, error)
	// Retire conditionally moves the active key to retiring and stamps
	// rotated_at. It reports whether the key was still active: false when a
	// concurrent rotation already retired it, in which case the caller must
	// not create a replacement (overlapping workers cannot double-create).
	Retire(ctx context.Context, id uuid.UUID, rotatedAt time.Time) (bool, error)
	// SweepRetiring retires the scope's keys whose rotated_at precedes
	// before: their grace period has elapsed, so nothing valid is signed by
	// them anymore.
	SweepRetiring(ctx context.Context, projectID *uuid.UUID, before time.Time) error
}
