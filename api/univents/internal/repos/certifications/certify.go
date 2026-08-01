package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) Certify(ctx context.Context, input models.CertifyInput) (*models.Certification, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.Certify")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).CreateCertification(ctx, sqlc.CreateCertificationParams{
		EditionID:        input.EditionID,
		TemplateID:       input.TemplateID,
		RegistrationID:   input.RegistrationID,
		UserID:           input.UserID,
		ProgramID:        input.ProgramID,
		VerificationHash: input.VerificationHash,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertification(result)), nil
}
