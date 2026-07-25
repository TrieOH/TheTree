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

var _ ports.ProfileRepo = (*Repo)(nil)

func NewProfileRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("actor_profile"),
	}
}

func mapActorProfile(src sqlc2.ActorProfile) models.ActorProfile {
	return models.ActorProfile{
		ActorID:   src.ActorID,
		Profile:   src.Profile,
		UpdatedAt: src.UpdatedAt,
	}
}
