package tos

import (
	"IdentityX/errs"
	"context"
	"time"

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

var _ ports.TosRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrTosNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapTos(src sqlc.TermsOfService) models.TermsOfService {
	return models.TermsOfService{
		ID:          src.ID,
		ProjectID:   src.ProjectID,
		Version:     src.Version,
		Content:     src.Content,
		EffectiveAt: src.EffectiveAt,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
}

func mapAcceptance(src sqlc.TosAcceptance) models.TosAcceptance {
	return models.TosAcceptance{
		ID:          src.ID,
		ActorID:     src.ActorID,
		ProjectID:   src.ProjectID,
		TosVersion:  src.TosVersion,
		ContentHash: src.ContentHash,
		Source:      src.Source,
		AcceptedAt:  src.AcceptedAt,
	}
}

func (repo *Repo) GetCurrent(ctx context.Context, projectID uuid.UUID) (*models.TermsOfService, error) {
	row, err := repo.db(ctx).GetCurrentTos(ctx, projectID)
	return database.One(row, err, repo.dbe, mapTos)
}

func (repo *Repo) Create(ctx context.Context, tos models.TermsOfService) (*models.TermsOfService, error) {
	row, err := repo.db(ctx).InsertTos(ctx, sqlc.InsertTosParams{
		ProjectID:   tos.ProjectID,
		Version:     tos.Version,
		Content:     tos.Content,
		EffectiveAt: tos.EffectiveAt,
	})
	return database.One(row, err, repo.dbe, mapTos)
}

func (repo *Repo) Update(ctx context.Context, projectID uuid.UUID, content string, effectiveAt time.Time) (*models.TermsOfService, error) {
	row, err := repo.db(ctx).UpdateTosBumpVersion(ctx, sqlc.UpdateTosBumpVersionParams{
		ProjectID:   projectID,
		Content:     content,
		EffectiveAt: effectiveAt,
	})
	return database.One(row, err, repo.dbe, mapTos)
}

func (repo *Repo) RecordAcceptance(ctx context.Context, acceptance models.TosAcceptance) (*models.TosAcceptance, error) {
	row, err := repo.db(ctx).InsertTosAcceptance(ctx, sqlc.InsertTosAcceptanceParams{
		ActorID:     acceptance.ActorID,
		ProjectID:   acceptance.ProjectID,
		TosVersion:  acceptance.TosVersion,
		ContentHash: acceptance.ContentHash,
		Source:      acceptance.Source,
	})
	return database.One(row, err, repo.dbe, mapAcceptance)
}

func (repo *Repo) LatestAcceptedVersion(ctx context.Context, actorID, projectID uuid.UUID) (int, error) {
	version, err := repo.db(ctx).GetLatestAcceptedTosVersion(ctx, sqlc.GetLatestAcceptedTosVersionParams{
		ActorID:   actorID,
		ProjectID: projectID,
	})
	return version, repo.dbe(err)
}

func (repo *Repo) ListAcceptances(ctx context.Context, projectID uuid.UUID) ([]models.TosAcceptanceWithActor, error) {
	rows, err := repo.db(ctx).ListTosAcceptances(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	out := make([]models.TosAcceptanceWithActor, 0, len(rows))
	for _, row := range rows {
		out = append(out, models.TosAcceptanceWithActor{
			ID:          row.ID,
			ActorID:     row.ActorID,
			ProjectID:   row.ProjectID,
			TosVersion:  row.TosVersion,
			ContentHash: row.ContentHash,
			Source:      row.Source,
			AcceptedAt:  row.AcceptedAt,
			ActorEmail:  row.ActorEmail,
		})
	}
	return out, nil
}
