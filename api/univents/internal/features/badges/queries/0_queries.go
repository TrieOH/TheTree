package queries

import "univents/ports"

type Queries struct {
	repo ports.BadgeTemplateRepo
}

func NewQueries(repo ports.BadgeTemplateRepo) *Queries {
	return &Queries{
		repo: repo,
	}
}
