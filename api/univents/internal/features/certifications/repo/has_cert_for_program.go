package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) HasCertForProgram(ctx context.Context, userID, programID uuid.UUID) (bool, error) {
	return database.Queries(ctx, repo.q).HasCertForProgram(ctx, sqlc.HasCertForProgramParams{
		UserID:    userID,
		ProgramID: &programID,
	})
}
