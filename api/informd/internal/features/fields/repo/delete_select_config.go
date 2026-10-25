package repo

import (
	"context"

	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) DeleteSelectConfig(ctx context.Context, fieldID uuid.UUID) error {
	err := database.Queries(ctx, repo.q).DeleteFieldSelectConfig(ctx, fieldID)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
