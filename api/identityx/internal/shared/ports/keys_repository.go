package ports

import (
	"context"

	"IdentityX/contracts"

	"github.com/google/uuid"
)

type KeysRepository interface {
	// --- Creation / Rotation ---

	CreateKeyPair(ctx context.Context, pair contracts.Pair) (*contracts.Pair, error)
	RotateSigningKeys(ctx context.Context, projectID *uuid.UUID) error

	// --- Signing (hot path) ---

	GetActiveSigningKey(ctx context.Context, projectID *uuid.UUID) (*contracts.Pair, error)
	GetActiveSigningKID(ctx context.Context, projectID *uuid.UUID) (string, error)

	// --- Verification ---

	GetKeyByKID(ctx context.Context, kid string) (*contracts.Pair, error)

	// --- Discovery (JWKS) ---

	ListPublicKeys(ctx context.Context, projectID *uuid.UUID) ([]contracts.PublicKey, error)

	// --- Revocation / Cleanup ---

	RevokeKeyByKID(ctx context.Context, kid string) error
	DeleteExpiredRevokedKeys(ctx context.Context) error
}
