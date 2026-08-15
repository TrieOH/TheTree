package crypto_keys

import (
	"context"
	"errors"
	"time"

	"IdentityX/internal/sqlc"
	"lib/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Retire conditionally moves the active key to retiring and stamps
// rotated_at. The WHERE status = 'active' guard makes the update
// conditional: when a concurrent rotation already retired the key, no row
// matches and Retire reports false, so the caller skips creating a
// replacement and overlapping workers cannot double-create.
func (repo *Repo) Retire(ctx context.Context, id uuid.UUID, rotatedAt time.Time) (bool, error) {
	_, err := database.Queries(ctx, repo.q).RetireCryptoKey(ctx, sqlc.RetireCryptoKeyParams{
		ID:        id,
		RotatedAt: &rotatedAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, repo.dbe(err)
	}
	return true, nil
}

// SweepRetiring retires the scope's keys whose grace period has elapsed:
// one RefreshTTL after rotation nothing valid is signed by them anymore.
func (repo *Repo) SweepRetiring(ctx context.Context, projectID *uuid.UUID, before time.Time) error {
	err := database.Queries(ctx, repo.q).SweepRetiringCryptoKeys(ctx, sqlc.SweepRetiringCryptoKeysParams{
		ProjectID: projectID,
		Before:    &before,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
