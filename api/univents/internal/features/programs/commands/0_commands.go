package commands

import (
	"univents/ports"
)

type Commands struct {
	events      ports.EventRepo
	editions    ports.EditionRepo
	programs    ports.ProgramRepo
	occurrences ports.ProgramOccurrenceRepo
}

func NewCommands(
	events ports.EventRepo,
	editions ports.EditionRepo,
	programs ports.ProgramRepo,
	occurrences ports.ProgramOccurrenceRepo,
) *Commands {
	return &Commands{
		events:      events,
		editions:    editions,
		programs:    programs,
		occurrences: occurrences,
	}
}
