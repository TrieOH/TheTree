package commands

import (
	"IdentityX/ports"
	"lib/database"
	"lib/errx"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewCommands(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return errx.MustProvide(&Commands{
		projects: projects,
		actors:   actors,
		logger:   logger,
		tracer:   tracer,
		tx:       tx,
	})
}
