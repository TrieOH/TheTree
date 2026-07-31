package repos

import (
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.WalletRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("wallet"),
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
