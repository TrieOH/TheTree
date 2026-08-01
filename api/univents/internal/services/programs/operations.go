package programs

import (
	"univents/ports"
)

type Operations struct {
	events      ports.EventRepo
	editions    ports.EditionRepo
	programs    ports.ProgramRepo
	occurrences ports.ProgramOccurrenceRepo
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	programs ports.ProgramRepo,
	occurrences ports.ProgramOccurrenceRepo,
) *Operations {
	return &Operations{
		events:      events,
		editions:    editions,
		programs:    programs,
		occurrences: occurrences,
	}
}
