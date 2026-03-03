package async

import (
	"univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type AsynqHandlers struct {
	activities domain.ActivitiesRepository
	gaClient   *goauth.Client
	tracer     trace.Tracer
	tx         database.TxRunner
}

func New(
	activities domain.ActivitiesRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		activities: activities,
		gaClient:   gaClient,
		tracer:     tracer,
		tx:         tx,
	}
}
