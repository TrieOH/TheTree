package feature_deps

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/ports"

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
