package commands

import (
	"univents/internal/eventcore/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	events   domain.EventsRepository
	gaClient *goauth.Client
	tracer   trace.Tracer
	tx       database.TxRunner
}

func New(
	events domain.EventsRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		events:   events,
		gaClient: gaClient,
		tracer:   tracer,
		tx:       tx,
	}
}
