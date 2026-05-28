package queries

import (
	"IdentityX/ports"
	"lib/database"
	"lib/errx"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	orgs   ports.OrganizationRepo
	logger *zap.Logger
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewQueries(deps ports.OrganizationDeps) *Queries {
	return errx.MustProvide(&Queries{
		orgs:   deps.Orgs,
		logger: deps.Logger,
		tracer: deps.Tracer,
		tx:     deps.Tx,
	})
}
