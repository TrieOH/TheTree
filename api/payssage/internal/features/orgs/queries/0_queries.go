package queries

import (
	"lib/database"
	"lib/errx"
	"payssage/ports"
	idx "sdk/identityx"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	orgs   ports.OrganizationRepo
	idx    *idx.Client
	logger *zap.Logger
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewQueries(
	orgs ports.OrganizationRepo,
	idx *idx.Client,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return errx.MustProvide(&Queries{
		orgs:   orgs,
		idx:    idx,
		logger: logger,
		tracer: tracer,
		tx:     tx,
	})
}
