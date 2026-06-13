package feature_deps

import (
	"time"

	"IdentityX/internal/shared/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ProjectCommandDeps struct {
	EncryptionKey []byte
	KeyLifetime   time.Duration
	Projects      ports.ProjectRepository
	Keys          ports.KeysRepository
	Logger        *zap.Logger
	Tracer        trace.Tracer
	Tx            database.TxRunner
}
