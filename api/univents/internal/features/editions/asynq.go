package editions

import (
	"context"
	"encoding/json"

	"lib/database"
	"univents/internal/platform/telemetry"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type AsynqHandlers struct {
	editions ports.EditionsRepository
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewAsynqService(
	editions ports.EditionsRepository,
	tracer trace.Tracer,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		editions: editions,
		tracer:   tracer,
		tx:       tx,
	}
}

func (uc *AsynqHandlers) HandleOpenEdition(ctx context.Context, t *asynq.Task) error {
	var payload contracts.EditionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		telemetry.Log().Error("HandleOpenEdition error", zap.Error(err))
		return err
	}

	return uc.editions.Open(ctx, payload.EditionID)
}

func (uc *AsynqHandlers) HandleStartEdition(ctx context.Context, t *asynq.Task) error {
	var payload contracts.EditionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		telemetry.Log().Error("HandleStartEdition error", zap.Error(err))
		return err
	}

	return uc.editions.Start(ctx, payload.EditionID)
}

func (uc *AsynqHandlers) HandleFinishEdition(ctx context.Context, t *asynq.Task) error {
	var payload contracts.EditionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		telemetry.Log().Error("HandleFinishEdition error", zap.Error(err))
		return err
	}

	return uc.editions.Finish(ctx, payload.EditionID)
}
