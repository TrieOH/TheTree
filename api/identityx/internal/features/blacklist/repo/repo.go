package repo

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"errors"
	"lib/database"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.BlacklistRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrBlacklistEntryNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

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

func (repo *Repo) Append(ctx context.Context, entry models.BlacklistEntry) error {
	_, err := repo.db(ctx).AppendBlacklistEntry(ctx, sqlc.AppendBlacklistEntryParams{
		CreatedByActorID: entry.CreatedByActorID,
		ProjectID:        entry.ProjectID,
		Type:             string(entry.Type),
		Target:           entry.Target,
		Reason:           entry.Reason,
		Metadata:         entry.Metadata,
		ExpiresAt:        entry.ExpiresAt,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return repo.dbe(err)
	}
	return nil
}

func (repo *Repo) GetByTarget(ctx context.Context, target string) (*models.BlacklistEntry, error) {
	sqlcEntry, err := repo.db(ctx).GetBlacklistEntryByTarget(ctx, target)
	return database.One(sqlcEntry, err, repo.dbe, mapEntry)
}

func (repo *Repo) GetByTargetAndType(ctx context.Context, target string, entryType models.BlacklistEntryType) (*models.BlacklistEntry, error) {
	sqlcEntry, err := repo.db(ctx).GetBlacklistEntryByTargetAndType(ctx, sqlc.GetBlacklistEntryByTargetAndTypeParams{
		Target: target,
		Type:   string(entryType),
	})
	return database.One(sqlcEntry, err, repo.dbe, mapEntry)
}
