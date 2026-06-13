package ports

import (
	"IdentityX/models"
	"context"
)

type BlacklistRepo interface {
	Append(ctx context.Context, entry models.BlacklistEntry) error
	GetByTarget(ctx context.Context, target string) (*models.BlacklistEntry, error)
	GetByTargetAndType(ctx context.Context, target string, entryType models.BlacklistEntryType) (*models.BlacklistEntry, error)
}
