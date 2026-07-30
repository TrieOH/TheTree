package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.ListByEdition")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListCertificationsByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapCertification), nil
}
