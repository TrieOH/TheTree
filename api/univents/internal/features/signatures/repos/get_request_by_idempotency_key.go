package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetRequestByIdempotencyKey(ctx context.Context, idempotencyKey uuid.UUID) (*models.SignatureRequest, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesRepo.GetRequestByIdempotencyKey")
	defer span.End()
	if err := repo.ExpireStaleRequests(ctx); err != nil {
		return nil, err
	}
	result, err := database.Queries(ctx, repo.q).GetSignatureRequestByIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignatureRequest(result)), nil
}
