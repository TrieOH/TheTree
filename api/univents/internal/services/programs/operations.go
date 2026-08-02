package programs

import (
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	events      ports.EventRepo
	editions    ports.EditionRepo
	programs    ports.ProgramRepo
	occurrences ports.ProgramOccurrenceRepo
	authz       *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	programs ports.ProgramRepo,
	occurrences ports.ProgramOccurrenceRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events:      events,
		editions:    editions,
		programs:    programs,
		occurrences: occurrences,
		authz:       authz,
	}
}
