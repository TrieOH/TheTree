package api_keys

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.APIKeysRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("api keys"),
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
