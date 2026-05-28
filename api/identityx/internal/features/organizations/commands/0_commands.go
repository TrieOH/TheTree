package commands

import (
	"IdentityX/ports"
	"lib/database"
	"lib/errx"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	actors ports.ActorRepo
	orgs   ports.OrganizationRepo
	logger *zap.Logger
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewCommands(deps ports.OrganizationDeps) *Commands {
	return errx.MustProvide(&Commands{
		actors: deps.Actors,
		orgs:   deps.Orgs,
		logger: deps.Logger,
		tracer: deps.Tracer,
		tx:     deps.Tx,
	})
}
