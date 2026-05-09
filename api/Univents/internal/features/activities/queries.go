package activities

import (
	"context"
	"univents/internal/platform/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	activities ports.ActivitiesRepository
	editions   ports.EditionsRepository
	tracer     trace.Tracer
	az         *authzed.Client
	tx         database.TxRunner
}

func NewQueryService(
	activities ports.ActivitiesRepository,
	editions ports.EditionsRepository,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		activities: activities,
		editions:   editions,
		tracer:     tracer,
		az:         az,
		tx:         tx,
	}
}

func (uc *QueryService) List(ctx context.Context, editionID uuid.UUID) (out []contracts.Activity, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ActivityService.List")
	defer span.End()

	return uc.activities.List(ctx, editionID)
}

func (uc *QueryService) AdminList(ctx context.Context, editionID uuid.UUID) (out []contracts.Activity, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ActivityService.AdminList")
	defer span.End()

	edition, err := uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_activities"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return nil, err
	}

	return uc.activities.ListAdmin(ctx, editionID)
}

func (uc *CommandService) ListRecords(ctx context.Context, activityID uuid.UUID) (records []contracts.AttendanceRecord, err error) {
	ctx, span := uc.tracer.Start(ctx, "ActivityService.ListRecords")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("mark.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var activity *contracts.Activity
	activity, err = uc.activities.GetByID(ctx, activityID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_attendance"),
		authz.Resource("activity", activity.ID.String()),
	); err != nil {
		return nil, err
	}

	var attendanceRecords []contracts.AttendanceRecord
	attendanceRecords, err = uc.activities.ListActivityAttendanceRecords(ctx, activityID)
	if err != nil {
		return nil, err
	}

	return attendanceRecords, nil
}
