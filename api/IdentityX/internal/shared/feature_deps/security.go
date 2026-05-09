package feature_deps

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type SecurityCommandDeps struct {
	Sessions ports.SessionRepository
	Project  ports.ProjectRepository
	Keys     ports.KeysRepository
	ApiKeys  ports.ApiKeyRepository
	Logger   *zap.Logger
	Tracer   trace.Tracer
	Tx       database.TxRunner
}
