package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.GetByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetCertificationByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertification(result)), nil
}
