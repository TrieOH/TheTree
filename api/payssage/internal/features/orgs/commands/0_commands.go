package commands

import (
	"lib/database"
	"lib/errx"
	"payssage/ports"
	idx "sdk/identityx"

	"go.opentelemetry.io/otel/trace"
)

type Commands struct {
	orgs   ports.OrganizationRepo
	idx    *idx.Client
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewCommands(
	orgs ports.OrganizationRepo,
	idx *idx.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return errx.MustProvide(&Commands{
		orgs:   orgs,
		idx:    idx,
		tracer: tracer,
		tx:     tx,
	})
}
