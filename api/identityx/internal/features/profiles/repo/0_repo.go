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

var _ ports.ProfileRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewNamedErrorHandler(errs.ErrProfileNotFound),
	}
}

func mapActorProfile(src sqlc.ActorProfile) models.ActorProfile {
	return models.ActorProfile{
		ActorID:       src.ActorID,
		Handle:        src.Handle,
		PfpURL:        src.PfpUrl,
		Profile:       src.Profile,
		SchemaVersion: src.SchemaVersion,
		Outdated:      src.Outdated,
		UpdatedAt:     src.UpdatedAt,
	}
}
