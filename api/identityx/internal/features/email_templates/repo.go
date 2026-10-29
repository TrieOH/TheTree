package email_templates

import (
	"IdentityX/errs"
	"context"

	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.EmailTemplateRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrEmailTemplateNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapEmailTemplate(src sqlc.EmailTemplate) models.EmailTemplate {
	return models.EmailTemplate{
		ID:        src.ID,
		ProjectID: src.ProjectID,
		Kind:      models.EmailTemplateKind(src.Kind),
		Subject:   src.Subject,
		Body:      src.Body,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func (repo *Repo) Upsert(ctx context.Context, template models.EmailTemplate) (*models.EmailTemplate, error) {
	sqlcTemplate, err := repo.db(ctx).UpsertEmailTemplate(ctx, sqlc.UpsertEmailTemplateParams{
		ProjectID: template.ProjectID,
		Kind:      string(template.Kind),
		Subject:   template.Subject,
		Body:      template.Body,
	})
	return database.One(sqlcTemplate, err, repo.dbe, mapEmailTemplate)
}

func (repo *Repo) GetByProjectAndKind(ctx context.Context, projectID uuid.UUID, kind models.EmailTemplateKind) (*models.EmailTemplate, error) {
	sqlcTemplate, err := repo.db(ctx).GetEmailTemplateByProjectAndKind(ctx, sqlc.GetEmailTemplateByProjectAndKindParams{
		ProjectID: projectID,
		Kind:      string(kind),
	})
	return database.One(sqlcTemplate, err, repo.dbe, mapEmailTemplate)
}

func (repo *Repo) Delete(ctx context.Context, projectID uuid.UUID, kind models.EmailTemplateKind) error {
	err := repo.db(ctx).DeleteEmailTemplate(ctx, sqlc.DeleteEmailTemplateParams{
		ProjectID: projectID,
		Kind:      string(kind),
	})
	return repo.dbe(err)
}

func (repo *Repo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.EmailTemplate, error) {
	rows, err := repo.db(ctx).ListEmailTemplatesByProject(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	out := make([]models.EmailTemplate, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapEmailTemplate(row))
	}
	return out, nil
}
