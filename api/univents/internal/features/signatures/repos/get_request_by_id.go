package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetRequestByID(ctx context.Context, id uuid.UUID) (*models.SignatureRequest, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesRepo.GetRequestByID")
	defer span.End()
	if err := repo.ExpireStaleRequests(ctx); err != nil {
		return nil, err
	}
	result, err := database.Queries(ctx, repo.q).GetSignatureRequestByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignatureRequest(result)), nil
}
