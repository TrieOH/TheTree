package repo

import (
	"context"
	"lib/database"
	"univents/models"
)

func (repo *Repo) GetByHash(ctx context.Context, hash string) (*models.Certification, error) {
	result, err := database.Queries(ctx, repo.q).GetCertificationByHash(ctx, hash)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertification(result)), nil
}
