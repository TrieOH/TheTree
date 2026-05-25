package ports

import (
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type AuthnDeps struct {
	Actors        ActorRepo
	PlatformRoles PlatformRolesRepo
	CryptoKeys    CryptoKeysRepo
	Blacklist     BlacklistRepo
	Logger        *zap.Logger
	Tracer        trace.Tracer
	Tx            database.TxRunner
}
