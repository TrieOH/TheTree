package repo

import (
	"lib/database"
	"payssage/errs"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.NamedErrorHandler
}

var _ ports.WalletRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewNamedErrorHandler(errs.ErrWalletNotFound),
	}
}

func mapWallet(src sqlc.Wallet) models.Wallet {
	return models.Wallet{
		ID:             src.ID,
		OwnerID:        src.OwnerID,
		OrganizationID: src.OrganizationID,
		Name:           src.Name,
		Sandbox:        src.Sandbox,
		FeeBps:         src.FeeBps,
		CollectorID:    src.CollectorID,
		CreatedAt:      src.CreatedAt,
	}
}
