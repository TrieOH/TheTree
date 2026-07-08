package activities

import (
	"context"
	"errors"
	idx "sdk/identityx"

	"lib/database"
	"univents/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"
	ports2 "univents/ports"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	activities ports.ActivitiesRepository
	editions   ports.EditionsRepository
	certs      ports2.CertificationRepo
	logger     *zap.Logger
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewCommandService(
	activities ports.ActivitiesRepository,
	editions ports.EditionsRepository,
	certs ports2.CertificationRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		activities: activities,
		editions:   editions,
		certs:      certs,
		logger:     logger,
		tracer:     tracer,
		tx:         tx,
	}
}

func (uc *CommandService) Create(ctx context.Context, in contracts.CreateActivitySpec) (out *contracts.Activity, err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, in.EditionID)
	if err != nil {
		return nil, err
	}

	var validActivity *contracts.Activity
	validActivity, err = contracts.NewActivity(ident.Sub.ID, in, edition)
	if err != nil {
		return nil, err
	}

	var created *contracts.Activity
	created, err = uc.activities.Create(ctx, validActivity)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *CommandService) Publish(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Publish")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("publish.success", err == nil))
	}()

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if activity.Status != contracts.ActivityStatusDraft {
		return errors.New("can't publish activities on statuses different than draft")
	}

	//var task *asynq.Task
	//task, err = contracts.NewStartActivityTask(activity.ID, activity.StartsAt)
	//if err != nil {
	//	return err
	//}
	//if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
	//	return err
	//}
	_ = uc.activities.Start(ctx, activity.ID)

	//task, err = contracts.NewEndActivityTask(activity.ID, activity.EndsAt)
	//if err != nil {
	//	return err
	//}
	//if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
	//	return err
	//}

	if err = uc.activities.Publish(ctx, activity.ID); err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) Register(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Register")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("register.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var isRegistered bool
	isRegistered, err = uc.activities.IsRegistered(ctx, ident.Sub.ID, activity.ID)
	if err != nil {
		return err
	}
	if isRegistered {
		return errx.Invalid("activity").SetMessage("user already registered to activity")
	}

	attendanceRecord := contracts.NewAttendanceRecord(ident.Sub.ID, activity.ID)
	if _, err = uc.activities.Register(ctx, *attendanceRecord); err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) Unregister(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.Unregister")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("unregister.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var isRegistered bool
	isRegistered, err = uc.activities.IsRegistered(ctx, ident.Sub.ID, activity.ID)
	if err != nil {
		return err
	}
	if !isRegistered {
		return errx.Invalid("activity").SetMessage("user isn't registered")
	}

	if err = uc.activities.Unregister(ctx, ident.Sub.ID, activity.ID); err != nil {
		return err
	}

	return nil
}

// FIXME Limit the times attendance can be marked to within the activity time

func (uc *CommandService) MarkAttendance(ctx context.Context, activityID, recordID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.MarkAttendance")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("mark.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	_, err = uc.activities.GetByID(ctx, activityID)
	if err != nil {
		return err
	}

	var attendanceRecord *contracts.AttendanceRecord
	attendanceRecord, err = uc.activities.GetAttendanceRecordByID(ctx, recordID)
	if err != nil {
		return err
	}
	if attendanceRecord.Status != contracts.AttendanceStatusRegistered {
		return errx.Invalid("attendance record").SetMessage("cannot mark attendance on activities on statuses different than registered")
	}

	if err = uc.activities.MarkAttendanceRecordStatus(ctx, recordID, &ident.Sub.ID, contracts.AttendanceStatusCompleted); err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) Complete(ctx context.Context, id uuid.UUID) (err error) {
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

		_, certErr := uc.certs.Certify(ctx, contracts.CertifyInput{
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
