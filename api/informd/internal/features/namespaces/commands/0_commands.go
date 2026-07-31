package commands

import (
	"Informd/ports"
)

type Commands struct {
	namespaces ports.NamespaceRepo
	forms      ports.FormsRepo
}

func NewCommands(
	projects ports.NamespaceRepo,
	forms ports.FormsRepo,
) *Commands {
	return &Commands{
		namespaces: projects,
		forms:      forms,
	}
}
