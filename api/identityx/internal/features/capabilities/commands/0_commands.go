package commands

import (
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	actors       ports.ActorRepo
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
	logger       *zap.Logger
	tracer       trace.Tracer
	tx           database.TxRunner
}

func NewCommands(
	actors ports.ActorRepo,
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		actors:       actors,
		capabilities: capabilities,
		projects:     projects,
		logger:       logger,
		tracer:       tracer,
		tx:           tx,
	}
}
