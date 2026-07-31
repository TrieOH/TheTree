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

var _ ports.CapabilityRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("capabilities"),
	}
}

func mapCapability(src sqlc.Capability) models.Capability {
	return models.Capability{
		ID:        src.ID,
		ProjectID: src.ProjectID,
		Resource:  src.Resource,
		Action:    src.Action,
		CreatedBy: src.CreatedBy,
		CreatedAt: src.CreatedAt,
	}
}
