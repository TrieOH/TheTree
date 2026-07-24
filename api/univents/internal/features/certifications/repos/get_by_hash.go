package repos

import (
	"context"
	"lib/database"
	"univents/models"
)

func (repo *repo) GetByHash(ctx context.Context, hash string) (*models.Certification, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByHash")
	defer span.End()
	cert, err := database.Queries(ctx, repo.q).GetCertificationByHash(ctx, &hash)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertification(cert)), nil
}
