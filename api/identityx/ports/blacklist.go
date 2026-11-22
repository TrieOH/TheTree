package ports

import (
	"IdentityX/models"
	"context"
)

type BlacklistRepo interface {
	// Append adds the entry and reports whether it did: false when the same
	// entry is already there.
	Append(ctx context.Context, entry models.BlacklistEntry) (bool, error)
	// Put adds the entry, or restarts it when it is already there: its
	// created_at moves to now.
	Put(ctx context.Context, entry models.BlacklistEntry) error
	GetByTarget(ctx context.Context, target string) (*models.BlacklistEntry, error)
	GetByTargetAndType(ctx context.Context, target string, entryType models.BlacklistEntryType) (*models.BlacklistEntry, error)
}
