package asynq

import (
	"univents/internal/commerce/domain"
	domain2 "univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type AsynqHandlers struct {
	tickets     domain.TicketsRepository
	products    domain.ProductsRepository
	activities  domain2.ActivitiesRepository
	checkpoints domain2.CheckpointsRepository
	gaClient    *goauth.Client
	tracer      trace.Tracer
	tx          database.TxRunner
}

func New(
	tickets domain.TicketsRepository,
	products domain.ProductsRepository,
	activities domain2.ActivitiesRepository,
	checkpoints domain2.CheckpointsRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		tickets:     tickets,
		products:    products,
		activities:  activities,
		checkpoints: checkpoints,
		gaClient:    gaClient,
		tracer:      tracer,
		tx:          tx,
	}
}
