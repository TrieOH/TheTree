package commands

import (
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	actors   ports.ActorRepo
	apiKeys  ports.ApiKeysRepo
	projects ports.ProjectRepo
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewCommands(
	actors ports.ActorRepo,
	apiKeys ports.ApiKeysRepo,
	projects ports.ProjectRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		actors:   actors,
		apiKeys:  apiKeys,
		projects: projects,
		logger:   logger,
		tracer:   tracer,
		tx:       tx,
	}
}
