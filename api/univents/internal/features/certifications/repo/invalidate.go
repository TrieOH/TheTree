package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) Invalidate(ctx context.Context, id uuid.UUID, reason *string) error {
	err := database.Queries(ctx, repo.q).InvalidateCertification(ctx, sqlc.InvalidateCertificationParams{
		ID:            id,
		InvalidReason: reason,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
