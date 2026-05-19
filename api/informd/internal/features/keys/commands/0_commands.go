package commands

import (
	"Informd/ports"
	"lib/authz"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	apiKeys  ports.ApiKeysRepo
	projects ports.NamespaceRepo
	perms    authz.Checker
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewCommands(
	apiKeys ports.ApiKeysRepo,
	projects ports.NamespaceRepo,
	perms authz.Checker,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		apiKeys:  apiKeys,
		projects: projects,
		perms:    perms,
		tx:       tx,
		tracer:   tracer,
	}
}
