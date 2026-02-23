package async

import (
	"univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type AsynqHandlers struct {
	editions domain.EditionsRepository
	gaClient *goauth.Client
	tracer   trace.Tracer
	tx       database.TxRunner
}

func New(
	editions domain.EditionsRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		editions: editions,
		gaClient: gaClient,
		tracer:   tracer,
		tx:       tx,
	}
}
