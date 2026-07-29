package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Signature, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesRepo.ListByEdition")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListSignaturesByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapSignature), nil
}
