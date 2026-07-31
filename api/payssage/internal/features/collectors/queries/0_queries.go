package queries

import (
	"payssage/ports"
)

type Queries struct {
	collectors ports.CollectorRepo
	orgs       ports.OrganizationRepo
}

func NewQueries(
	collectors ports.CollectorRepo,
	orgs ports.OrganizationRepo,
) *Queries {
	return &Queries{
		collectors: collectors,
		orgs:       orgs,
	}
}
