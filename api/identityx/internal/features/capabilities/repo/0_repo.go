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

var _ ports.CapabilityRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewNamedErrorHandler(errs.ErrCapabilityNotFound),
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
