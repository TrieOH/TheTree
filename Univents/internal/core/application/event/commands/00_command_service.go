package commands

import (
	"univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	events   domain.EventsRepository
	minio    *minio.Client
	gaClient *goauth.Client
	tracer   trace.Tracer
	tx       database.TxRunner
}

func New(
	events domain.EventsRepository,
	minio *minio.Client,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		events:   events,
		minio:    minio,
		gaClient: gaClient,
		tracer:   tracer,
		tx:       tx,
	}
}
