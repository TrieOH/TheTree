package outbounds

import (
	"GoAuth/internal/domain/key"
	"context"

	"github.com/google/uuid"
)

type KeysRepository interface {
	// --- Creation / Rotation ---

	CreateKeyPair(ctx context.Context, pair key.Pair) (*key.Pair, error)
	RotateGoAuthSigningKeys(ctx context.Context) error
	RotateProjectSigningKeys(ctx context.Context, projectID uuid.UUID) error

	// --- Signing (hot path) ---

	GetActiveGoAuthSigningKey(ctx context.Context) (*key.Pair, error)
	GetActiveProjectSigningKey(ctx context.Context, projectID uuid.UUID) (*key.Pair, error)

	// --- Verification ---

	GetGoAuthKeyByKID(ctx context.Context, kid string) (*key.Pair, error)
	GetProjectKeyByKID(ctx context.Context, kid string) (*key.Pair, error)

	// --- Discovery (JWKS) ---

	ListGoAuthPublicKeys(ctx context.Context) ([]key.PublicKey, error)
	ListProjectPublicKeys(ctx context.Context, projectID uuid.UUID) ([]key.PublicKey, error)

	// --- Revocation / Cleanup ---

	RevokeKeyByKID(ctx context.Context, kid string) error
	DeleteExpiredRevokedKeys(ctx context.Context) error

	// --- Metadata ---

	GetActiveGoAuthSigningKID(ctx context.Context) (string, error)
	GetActiveProjectSigningKID(ctx context.Context, projectID uuid.UUID) (string, error)
}
