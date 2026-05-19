package commands

import (
	"Informd/ports"
	"lib/authz"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	namespaces ports.NamespaceRepo
	perms      authz.Checker
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	projects ports.NamespaceRepo,
	perms authz.Checker,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		namespaces: projects,
		perms:      perms,
		tx:         tx,
		tracer:     tracer,
	}
}
