package commands

import (
	"univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	events   domain.EventsRepository
	minio    *minio.Client
	gaClient *goauth.Client
	tracer   trace.Tracer
	az       *authzed.Client
	tx       database.TxRunner
}

func New(
	events domain.EventsRepository,
	minio *minio.Client,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		events:   events,
		minio:    minio,
		gaClient: gaClient,
		tracer:   tracer,
		az:       az,
		tx:       tx,
	}
}
