package queries

import (
	"IdentityX/ports"
	"lib/database"
	"lib/errx"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	projects ports.ProjectRepo
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewQueries(
	projects ports.ProjectRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return errx.MustProvide(&Queries{
		projects: projects,
		logger:   logger,
		tracer:   tracer,
		tx:       tx,
	})
}
