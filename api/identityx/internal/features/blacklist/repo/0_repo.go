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

var _ ports.BlacklistRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewNamedErrorHandler(errs.ErrBlacklistEntryNotFound),
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
