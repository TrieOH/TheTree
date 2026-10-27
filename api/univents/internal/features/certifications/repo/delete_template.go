package repo

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	err := database.Queries(ctx, repo.q).DeleteCertificationTemplate(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
