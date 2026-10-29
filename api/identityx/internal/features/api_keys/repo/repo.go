package repo

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.APIKeysRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrAPIKeyNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapAPIKey(src sqlc.ApiKey) models.APIKey {
	return models.APIKey{
		ID:            src.ID,
		SubjectID:     src.SubjectID,
		Name:          src.Name,
		DisplayPrefix: src.DisplayPrefix,
		KeyHash:       src.KeyHash,
		Metadata:      src.Metadata,
		ExpiresAt:     src.ExpiresAt,
		RevokedAt:     src.RevokedAt,
		LastUsedAt:    src.LastUsedAt,
		CreatedBy:     src.CreatedBy,
		CreatedAt:     src.CreatedAt,
	}
}

func (repo *Repo) Create(ctx context.Context, toCreate models.APIKey) (*models.APIKey, error) {
	row, err := repo.db(ctx).CreateApiKey(ctx, sqlc.CreateApiKeyParams{
		SubjectID:     toCreate.SubjectID,
		Name:          toCreate.Name,
		DisplayPrefix: toCreate.DisplayPrefix,
		KeyHash:       toCreate.KeyHash,
		CreatedBy:     toCreate.CreatedBy,
		ExpiresAt:     toCreate.ExpiresAt,
	})
	return database.One(row, err, repo.dbe, mapAPIKey)
}

func (repo *Repo) GetByPrefix(ctx context.Context, prefix string) (*models.APIKey, error) {
	row, err := repo.db(ctx).GetApiKeyByPrefix(ctx, prefix)
	return database.One(row, err, repo.dbe, mapAPIKey)
}
