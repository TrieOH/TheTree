package ports

import (
	"IdentityX/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type KeysRepository interface {
	// --- Creation / Rotation ---

	CreateKeyPair(ctx context.Context, pair contracts.Pair) (*contracts.Pair, error)
	RotateGoAuthSigningKeys(ctx context.Context) error
	RotateProjectSigningKeys(ctx context.Context, projectID uuid.UUID) error

	// --- Signing (hot path) ---

	GetActiveGoAuthSigningKey(ctx context.Context) (*contracts.Pair, error)
	GetActiveProjectSigningKey(ctx context.Context, projectID uuid.UUID) (*contracts.Pair, error)

	// --- Verification ---

	GetGoAuthKeyByKID(ctx context.Context, kid string) (*contracts.Pair, error)
	GetProjectKeyByKID(ctx context.Context, kid string) (*contracts.Pair, error)

	// --- Discovery (JWKS) ---

	ListGoAuthPublicKeys(ctx context.Context) ([]contracts.PublicKey, error)
	ListProjectPublicKeys(ctx context.Context, projectID uuid.UUID) ([]contracts.PublicKey, error)

	// --- Revocation / Cleanup ---

	RevokeKeyByKID(ctx context.Context, kid string) error
	DeleteExpiredRevokedKeys(ctx context.Context) error

	// --- Metadata ---

	GetActiveGoAuthSigningKID(ctx context.Context) (string, error)
	GetActiveProjectSigningKID(ctx context.Context, projectID uuid.UUID) (string, error)
}
