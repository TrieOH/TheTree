package queries

import (
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	profiles ports.ProfileRepo
	projects ports.ProjectRepo
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func New(
	profiles ports.ProfileRepo,
	projects ports.ProjectRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		profiles: profiles,
		projects: projects,
		logger:   logger,
		tracer:   tracer,
		tx:       tx,
	}
}
