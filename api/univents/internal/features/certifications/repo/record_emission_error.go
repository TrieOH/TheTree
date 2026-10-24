package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) RecordEmissionError(ctx context.Context, input *models.CertEmissionError) error {
	err := database.Queries(ctx, repo.q).RecordCertEmissionError(ctx, sqlc.RecordCertEmissionErrorParams{
		EditionID:    input.EditionID,
		UserID:       input.UserID,
		TemplateID:   input.TemplateID,
		ProgramID:    input.ProgramID,
		ErrorMessage: input.ErrorMessage,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
