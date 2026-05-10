package feature_deps

import (
	"IdentityX/internal/shared/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type SessionCommandDeps struct {
	Sessions ports.SessionRepository
	Keys     ports.KeysRepository
	Logger   *zap.Logger
	Tracer   trace.Tracer
	Tx       database.TxRunner
}
