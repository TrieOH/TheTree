package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.DeleteTemplate")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteCertificationTemplate(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
