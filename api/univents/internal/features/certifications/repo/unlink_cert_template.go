package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) UnlinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	return database.Queries(ctx, repo.q).DeleteCertificationTemplatePrograms(ctx, sqlc.DeleteCertificationTemplateProgramsParams{
		TemplateID: templateID,
		ProgramID:  programID,
	})
}
