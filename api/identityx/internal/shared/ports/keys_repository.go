package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type KeysRepository interface {
	// --- Creation / Rotation ---

	CreateKeyPair(ctx context.Context, pair models.Pair) (*models.Pair, error)
	RotateSigningKeys(ctx context.Context, projectID *uuid.UUID) error

	// --- Signing (hot path) ---

	GetActiveSigningKey(ctx context.Context, projectID *uuid.UUID) (*models.Pair, error)
	GetActiveSigningKID(ctx context.Context, projectID *uuid.UUID) (string, error)

	// --- Verification ---

	GetKeyByKID(ctx context.Context, kid string) (*models.Pair, error)

	// --- Discovery (JWKS) ---

	ListPublicKeys(ctx context.Context, projectID *uuid.UUID) ([]models.PublicKey, error)

	// --- Revocation / Cleanup ---

	RevokeKeyByKID(ctx context.Context, kid string) error
	DeleteExpiredRevokedKeys(ctx context.Context) error
}
