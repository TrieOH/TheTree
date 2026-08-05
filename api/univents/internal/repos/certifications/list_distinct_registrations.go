package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListDistinctRegistrationsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.CertEligibleAttendee, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.ListDistinctRegistrationsByEdition")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListDistinctRegistrationsByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapEligibleAttendee), nil
}
