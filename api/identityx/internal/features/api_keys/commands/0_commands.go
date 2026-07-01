package commands

import (
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	hmacSecret   []byte
	actors       ports.ActorRepo
	apiKeys      ports.ApiKeysRepo
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
	logger       *zap.Logger
	tracer       trace.Tracer
	tx           database.TxRunner
}

func NewCommands(
	hmacSecret []byte,
	actors ports.ActorRepo,
	apiKeys ports.ApiKeysRepo,
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		hmacSecret:   hmacSecret,
		actors:       actors,
		apiKeys:      apiKeys,
		capabilities: capabilities,
		projects:     projects,
		logger:       logger,
		tracer:       tracer,
		tx:           tx,
	}
}
