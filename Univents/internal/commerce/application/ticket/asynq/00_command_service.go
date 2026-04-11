package asynq

import (
	"univents/internal/commerce/domain"
	domain2 "univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"go.opentelemetry.io/otel/trace"
)

type AsynqHandlers struct {
	tickets     domain.TicketsRepository
	products    domain.ProductsRepository
	activities  domain2.ActivitiesRepository
	checkpoints domain2.CheckpointsRepository
	gaClient    *goauth.Client
	tracer      trace.Tracer
	az          *authzed.Client
	tx          database.TxRunner
}

func New(
	tickets domain.TicketsRepository,
	products domain.ProductsRepository,
	activities domain2.ActivitiesRepository,
	checkpoints domain2.CheckpointsRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		tickets:     tickets,
		products:    products,
		activities:  activities,
		checkpoints: checkpoints,
		gaClient:    gaClient,
		tracer:      tracer,
		az:          az,
		tx:          tx,
	}
}
