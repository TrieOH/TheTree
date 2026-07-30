package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) HasCertForProgram(ctx context.Context, userID, programID uuid.UUID) (bool, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.HasCertForProgram")
	defer span.End()
	return database.Queries(ctx, repo.q).HasCertForProgram(ctx, sqlc.HasCertForProgramParams{
		UserID:    userID,
		ProgramID: &programID,
	})
}
