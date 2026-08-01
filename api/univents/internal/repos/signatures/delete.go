package signatures

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesRepo.Delete")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteSignature(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
