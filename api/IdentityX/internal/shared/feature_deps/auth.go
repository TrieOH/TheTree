package feature_deps

import (
	"IdentityX/internal/shared/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type AuthCommandDeps struct {
	EncryptionKey []byte
	Issuer        string
	Users         ports.UserRepository
	Sessions      ports.SessionRepository
	Projects      ports.ProjectRepository
	Keys          ports.KeysRepository
	Renderer      ports.EmailRenderer
	Mailer        ports.Mailer
	Logger        *zap.Logger
	Tracer        trace.Tracer
	Tx            database.TxRunner
}
