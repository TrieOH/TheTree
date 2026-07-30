package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListDistinctParticipantsByProgram(ctx context.Context, programID uuid.UUID) ([]models.CertEligibleAttendee, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.ListDistinctParticipantsByProgram")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListDistinctParticipantsByProgram(ctx, programID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapEligibleParticipant), nil
}
