package activities

import (
	"context"
	"errors"

	"lib/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	activities ports.ActivitiesRepository
	editions   ports.EditionsRepository
	logger     *zap.Logger
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewCommandService(
	activities ports.ActivitiesRepository,
	editions ports.EditionsRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		activities: activities,
		editions:   editions,
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

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, in.EditionID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_activity"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return nil, err
	}

	var validActivity *contracts.Activity
	validActivity, err = contracts.NewActivity(sub.ID, in, edition)
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

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("publish"),
		authz.Resource("activity", activity.ID.String()),
	); err != nil {
		return err
	}

	if activity.Status != contracts.ActivityStatusDraft {
		return errors.New("can't publish activities on statuses different than draft")
	}

	var task *asynq.Task
	task, err = contracts.NewStartActivityTask(activity.ID, activity.StartsAt)
	if err != nil {
		return err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		return err
	}

	task, err = contracts.NewEndActivityTask(activity.ID, activity.EndsAt)
	if err != nil {
		return err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		return err
	}

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

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var isRegistered bool
	isRegistered, err = uc.activities.IsRegistered(ctx, sub.ID, activity.ID)
	if err != nil {
		return err
	}
	if isRegistered {
		return errx.Invalid("activity").SetMessage("user already registered to activity")
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("attend"),
		authz.Resource("activity", activity.ID.String()),
	); err != nil {
		return err
	}

	attendanceRecord := contracts.NewAttendanceRecord(sub.ID, activity.ID)
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

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var isRegistered bool
	isRegistered, err = uc.activities.IsRegistered(ctx, sub.ID, activity.ID)
	if err != nil {
		return err
	}
	if !isRegistered {
		return errx.Invalid("activity").SetMessage("user isn't registered")
	}

	if err = uc.activities.Unregister(ctx, sub.ID, activity.ID); err != nil {
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

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, activityID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("mark_attendance"),
		authz.Resource("activity", activity.ID.String()),
	); err != nil {
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

	if err = uc.activities.MarkAttendanceRecordStatus(ctx, recordID, &sub.ID, contracts.AttendanceStatusCompleted); err != nil {
		return err
	}

	return nil
}
