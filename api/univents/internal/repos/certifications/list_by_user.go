package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.ListByUser")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListCertificationsByUser(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapCertification), nil
}
