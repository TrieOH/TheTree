package queries

import (
	"IdentityX/ports"
	"lib/errx"
)

type Queries struct {
	projects ports.ProjectRepo
}

func NewQueries(
	projects ports.ProjectRepo,
) *Queries {
	return errx.MustProvide(&Queries{
		projects: projects,
	})
}
