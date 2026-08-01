package collectors

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

var _ ports.CollectorRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("collector"),
	}
}

func mapCollector(src sqlc.Collector) models.Collector {
	return models.Collector{
		ID:             src.ID,
		OwnerID:        src.OwnerID,
		OrganizationID: src.OrganizationID,
		Provider:       src.Provider,
		ProviderUserID: src.ProviderUserID,
		Credentials:    src.Credentials,
		CreatedAt:      src.CreatedAt,
		RevokedAt:      src.RevokedAt,
	}
}
