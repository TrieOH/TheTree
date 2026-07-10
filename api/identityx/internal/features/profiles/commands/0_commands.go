package commands

import (
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	profiles ports.ProfileRepo
	schemas  ports.ProfileSchemaRepo
	projects ports.ProjectRepo
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func New(
	profiles ports.ProfileRepo,
	schemas ports.ProfileSchemaRepo,
	projects ports.ProjectRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		profiles: profiles,
		schemas:  schemas,
		projects: projects,
		logger:   logger,
		tracer:   tracer,
		tx:       tx,
	}
}
