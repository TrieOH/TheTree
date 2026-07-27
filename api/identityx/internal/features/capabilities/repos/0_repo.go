package repos

import (
	sqlc2 "IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc2.Queries
	dbe database.ErrorHandler
}

var _ ports.CapabilityRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("capabilities"),
	}
}

func mapCapability(src sqlc2.Capability) models.Capability {
	return models.Capability{
		ID:        src.ID,
		ProjectID: src.ProjectID,
		Resource:  src.Resource,
		Action:    src.Action,
		CreatedBy: src.CreatedBy,
		CreatedAt: src.CreatedAt,
	}
}
