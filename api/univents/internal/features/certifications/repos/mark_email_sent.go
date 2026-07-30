package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) MarkEmailSent(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.MarkEmailSent")
	defer span.End()
	err := database.Queries(ctx, repo.q).MarkCertificationEmailSent(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
