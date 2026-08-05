package collectors

import (
	"payssage/internal/authz"

	"payssage/ports"
)

type Operations struct {
	collectors ports.CollectorRepo
	orgs       ports.OrganizationRepo
	authz      *authz.Service
}

func NewOperations(
	collectors ports.CollectorRepo,
	orgs ports.OrganizationRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		collectors: collectors,
		orgs:       orgs,
		authz:      authz,
	}
}
