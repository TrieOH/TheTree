package editions

import (
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	authz    *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events:   events,
		editions: editions,
		authz:    authz,
	}
}
