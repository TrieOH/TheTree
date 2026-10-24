package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"errors"
	"lib/database"

	"github.com/jackc/pgx/v5"
)

func (repo *Repo) Append(ctx context.Context, entry models.BlacklistEntry) error {
	_, err := database.Queries(ctx, repo.q).AppendBlacklistEntry(ctx, sqlc.AppendBlacklistEntryParams{
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
