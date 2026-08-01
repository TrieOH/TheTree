package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListEmissionErrorsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.CertEmissionError, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.ListEmissionErrorsByEdition")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListCertEmissionErrorsByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapCertEmissionError), nil
}
