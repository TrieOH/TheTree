package queries

import (
	"univents/ports"
)

type Queries struct {
	programs    ports.ProgramRepo
	occurrences ports.ProgramOccurrenceRepo
}

func NewQueries(
	programs ports.ProgramRepo,
	occurrences ports.ProgramOccurrenceRepo,
) *Queries {
	return &Queries{
		programs:    programs,
		occurrences: occurrences,
	}
}
