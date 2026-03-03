package async

import (
	"context"
	"encoding/json"
	"univents/internal/core/domain"
	"univents/internal/plataform/telemetry"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func (uc *AsynqHandlers) HandleOpenEdition(ctx context.Context, t *asynq.Task) error {
	var payload domain.EditionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		telemetry.Log().Error("HandleOpenEdition error", zap.Error(err))
		return err
	}

	return uc.editions.Open(ctx, payload.EditionID)
}

func (uc *AsynqHandlers) HandleStartEdition(ctx context.Context, t *asynq.Task) error {
	var payload domain.EditionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		telemetry.Log().Error("HandleStartEdition error", zap.Error(err))
		return err
	}

	return uc.editions.Start(ctx, payload.EditionID)
}

func (uc *AsynqHandlers) HandleFinishEdition(ctx context.Context, t *asynq.Task) error {
	var payload domain.EditionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		telemetry.Log().Error("HandleFinishEdition error", zap.Error(err))
		return err
	}

	return uc.editions.Finish(ctx, payload.EditionID)
}
