package inbounds

import (
	"context"

	"github.com/google/uuid"
)

type KeysService interface {
	// --- Signing ---

	SignGoAuth(ctx context.Context, payload []byte) (sig []byte, err error)
	SignProject(ctx context.Context, projectID uuid.UUID, payload []byte) (sig []byte, err error)

	// --- Verification ---

	VerifyGoAuth(ctx context.Context, kid string, payload, sig []byte) error
	VerifyProject(ctx context.Context, projectID uuid.UUID, kid string, payload, sig []byte) error

	// --- Metadata ---

	GetActiveGoAuthSigningKID(ctx context.Context) (string, error)
	GetActiveProjectSigningKID(ctx context.Context, projectID uuid.UUID) (string, error)

	// --- Management ---

	RevokeKey(ctx context.Context, kid string) error
}
