package editions

import (
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	badges   ports.BadgeEditionOps
	authz    *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	authz *authz.Service,
	badges ports.BadgeEditionOps,
) *Operations {
	return &Operations{
		events:   events,
		editions: editions,
		badges:   badges,
		authz:    authz,
	}
}
