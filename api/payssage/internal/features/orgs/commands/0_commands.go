package commands

import (
	"lib/database"
	"lib/errx"
	"payssage/ports"
	idx "sdk/identityx"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	orgs   ports.OrganizationRepo
	idx    *idx.Client
	logger *zap.Logger
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewCommands(
	orgs ports.OrganizationRepo,
	idx *idx.Client,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return errx.MustProvide(&Commands{
		orgs:   orgs,
		idx:    idx,
		logger: logger,
		tracer: tracer,
		tx:     tx,
	})
}
