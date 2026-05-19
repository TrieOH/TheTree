package commands

import (
	"Informd/ports"
	"lib/authz"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	namespaces ports.NamespaceRepo
	perms      authz.Checker
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
	perms authz.Checker,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
		perms:      perms,
		tx:         tx,
		tracer:     tracer,
	}
}
