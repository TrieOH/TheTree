package signatures

import (
	"context"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) ExpireStaleRequests(ctx context.Context) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesRepo.ExpireStaleRequests")
	defer span.End()
	err := database.Queries(ctx, repo.q).ExpireStaleSignatureRequests(ctx)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
