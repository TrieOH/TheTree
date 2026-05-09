package feature_deps

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/ports"
	"time"

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
