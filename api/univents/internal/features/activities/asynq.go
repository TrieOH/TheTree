package activities

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
	activities ports.ActivitiesRepository
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewAsynqService(
	activities ports.ActivitiesRepository,
	tracer trace.Tracer,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		activities: activities,
		tracer:     tracer,
		tx:         tx,
	}
}

func (uc *AsynqHandlers) HandleStartActivity(ctx context.Context, t *asynq.Task) error {
	var payload contracts.ActivityPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		telemetry.Log().Error("HandleStartActivity error", zap.Error(err))
		return err
	}

	return uc.activities.Start(ctx, payload.ActivityID)
}

func (uc *AsynqHandlers) HandleFinishActivity(ctx context.Context, t *asynq.Task) error {
	var payload contracts.ActivityPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		telemetry.Log().Error("HandleFinishActivity error", zap.Error(err))
		return err
	}

	return uc.activities.Finish(ctx, payload.ActivityID)
}
