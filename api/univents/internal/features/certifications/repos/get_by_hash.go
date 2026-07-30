package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"
)

func (repo *Repo) GetByHash(ctx context.Context, hash string) (*models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.GetByHash")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetCertificationByHash(ctx, hash)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertification(result)), nil
}
