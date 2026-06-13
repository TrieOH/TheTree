package ports

import (
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type OrganizationDeps struct {
	Actors ActorRepo
	Orgs   OrganizationRepo
	Logger *zap.Logger
	Tracer trace.Tracer
	Tx     database.TxRunner
}
