package collectors

import (
	"payssage/ports"
)

type Operations struct {
	collectors ports.CollectorRepo
	orgs       ports.OrganizationRepo
}

func NewOperations(
	collectors ports.CollectorRepo,
	orgs ports.OrganizationRepo,
) *Operations {
	return &Operations{
		collectors: collectors,
		orgs:       orgs,
	}
}
