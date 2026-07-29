package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Signature, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesRepo.GetByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetSignatureByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignature(result)), nil
}
