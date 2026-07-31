package repos

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

var _ ports.BlacklistRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("blacklist entry"),
	}
}

func mapEntry(src sqlc.BlacklistEntry) models.BlacklistEntry {
	return models.BlacklistEntry{
		ID:               src.ID,
		CreatedByActorID: src.CreatedByActorID,
		ProjectID:        src.ProjectID,
		Type:             models.BlacklistEntryType(src.Type),
		Target:           src.Target,
		Reason:           src.Reason,
		Metadata:         src.Metadata,
		CreatedAt:        src.CreatedAt,
		ExpiresAt:        src.ExpiresAt,
	}
}
