package commands

import (
	"context"
	"univents/contracts"
	"univents/models"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (uc *Commands) Complete(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Complete")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("complete.success", err == nil))
	}()

	if err = uc.activities.Finish(ctx, id); err != nil {
		return err
	}

	records, err := uc.activities.ListActivityAttendanceRecords(ctx, id)
	if err != nil {
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
			TargetID:   id,
			TargetType: "activity",
		})
		if certErr != nil {
			uc.logger.Warn("failed to certify user on activity complete",
				zap.String("user_id", record.UserID.String()),
				zap.String("activity_id", id.String()),
				zap.Error(certErr),
			)
		}
	}

	return nil
}
