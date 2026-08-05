package profiles

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

var _ ports.ProfileRepo = (*Repo)(nil)

func NewProfileRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("actor_profile"),
	}
}

func mapActorProfile(src sqlc.ActorProfile) models.ActorProfile {
	return models.ActorProfile{
		ActorID:       src.ActorID,
		Handle:        src.Handle,
		Profile:       src.Profile,
		SchemaVersion: src.SchemaVersion,
		Outdated:      src.Outdated,
		UpdatedAt:     src.UpdatedAt,
	}
}
