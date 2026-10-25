package repo

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.NamedErrorHandler
}

var _ ports.APIKeysRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewNamedErrorHandler(errs.ErrAPIKeyNotFound),
	}
}

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
