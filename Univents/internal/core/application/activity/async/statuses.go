package async

import (
	"context"
	"encoding/json"
	"univents/internal/core/domain"
	"univents/internal/plataform/telemetry"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func (uc *AsynqHandlers) HandleStartActivity(ctx context.Context, t *asynq.Task) error {
	var payload domain.ActivityPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		telemetry.Log().Error("HandleStartActivity error", zap.Error(err))
		return err
	}

	return uc.activities.Start(ctx, payload.ActivityID)
}

func (uc *AsynqHandlers) HandleFinishActivity(ctx context.Context, t *asynq.Task) error {
	var payload domain.ActivityPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		telemetry.Log().Error("HandleFinishActivity error", zap.Error(err))
		return err
	}

	return uc.activities.Finish(ctx, payload.ActivityID)
}
