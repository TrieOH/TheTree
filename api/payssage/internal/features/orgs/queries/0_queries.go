package queries

import (
	"lib/database"
	"lib/errx"
	"payssage/ports"
	idx "sdk/identityx"

	"go.opentelemetry.io/otel/trace"
)

type Queries struct {
	orgs   ports.OrganizationRepo
	idx    *idx.Client
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewQueries(
	orgs ports.OrganizationRepo,
	idx *idx.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return errx.MustProvide(&Queries{
		orgs:   orgs,
		idx:    idx,
		tracer: tracer,
		tx:     tx,
	})
}
