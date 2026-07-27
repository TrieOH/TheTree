package queries

import (
	"IdentityX/ports"
)

type Queries struct {
	projects ports.ProjectRepo
	actors   ports.ActorRepo
}

func NewQueries(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
) *Queries {
	return &Queries{
		projects: projects,
		actors:   actors,
	}
}
