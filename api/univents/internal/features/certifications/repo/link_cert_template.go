package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) LinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	return database.Queries(ctx, repo.q).CreateCertificationTemplateProgram(ctx, sqlc.CreateCertificationTemplateProgramParams{
		TemplateID: templateID,
		ProgramID:  programID,
	})
}
