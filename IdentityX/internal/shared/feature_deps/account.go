package feature_deps

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type AccountCommandDeps struct {
	Users          ports.UserRepository
	Accounts       ports.AccountRepository
	Sessions       ports.SessionRepository
	Keys           ports.KeysRepository
	TokenReuseList ports.TokenReuseListRepository
	MailRenderer   ports.EmailRenderer
	MailSender     ports.Mailer
	Logger         *zap.Logger
	Tracer         trace.Tracer
	Tx             database.TxRunner
}
