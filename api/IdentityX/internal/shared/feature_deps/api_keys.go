package feature_deps

import (
	"IdentityX/internal/shared/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ApiKeysCommandDeps struct {
	ApiKeys ports.ApiKeyRepository
	Project ports.ProjectRepository
	Logger  *zap.Logger
	Tracer  trace.Tracer
	Tx      database.TxRunner
}
