package activities

import (
	"context"
	"encoding/json"
	"univents/models"

	"lib/database"
	"lib/telemetry"
	"univents/contracts"
	"univents/internal/shared/ports"
	ports2 "univents/ports"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type AsynqHandlers struct {
	activities ports.ActivitiesRepository
	certs      ports2.CertificationRepo
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewAsynqService(
	activities ports.ActivitiesRepository,
	certs ports2.CertificationRepo,
	tracer trace.Tracer,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		activities: activities,
		certs:      certs,
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

	if err := uc.activities.Finish(ctx, payload.ActivityID); err != nil {
		return err
	}

	records, err := uc.activities.ListActivityAttendanceRecords(ctx, payload.ActivityID)
	if err != nil {
		telemetry.Log().Error("HandleFinishActivity: list records error", zap.Error(err))
		return err
	}

	certified := make(map[uuid.UUID]bool)
	for _, record := range records {
		if record.Status != contracts.AttendanceStatusCompleted {
			continue
		}
		if certified[record.UserID] {
			continue
		}
		certified[record.UserID] = true

		_, certErr := uc.certs.Certify(ctx, models.CertifyInput{
			UserID:     record.UserID,
			TargetID:   payload.ActivityID,
			TargetType: "activity",
		})
		if certErr != nil {
			telemetry.Log().Warn("HandleFinishActivity: certify error",
				zap.String("user_id", record.UserID.String()),
				zap.String("activity_id", payload.ActivityID.String()),
				zap.Error(certErr),
			)
		}
	}

	return nil
}
