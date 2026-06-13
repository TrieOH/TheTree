package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"univents/internal/shared/errx"
	"univents/internal/shared/validation"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type ActivityStatus string

const (
	ActivityStatusDraft     ActivityStatus = "draft"
	ActivityStatusPublished ActivityStatus = "published"
	ActivityStatusOngoing   ActivityStatus = "ongoing"
	ActivityStatusCompleted ActivityStatus = "completed"
	ActivityStatusCanceled  ActivityStatus = "canceled"
)

type DifficultyLevel string

const (
	DifficultyLevelNoPrerequisites DifficultyLevel = "no_prerequisites"
	DifficultyLevelBeginner        DifficultyLevel = "beginner"
	DifficultyLevelIntermediate    DifficultyLevel = "intermediate"
	DifficultyLevelAdvanced        DifficultyLevel = "advanced"
	DifficultyLevelExpert          DifficultyLevel = "expert"
)

type Activity struct {
	ID                uuid.UUID        `json:"id"`
	EditionID         uuid.UUID        `json:"edition_id" validate:"required"`
	Title             string           `json:"title"      validate:"required,min=3"`
	Description       *string          `json:"description"`
	Status            ActivityStatus   `json:"status"     validate:"required,oneof=draft published ongoing completed canceled"`
	Location          string           `json:"location"`
	StartsAt          time.Time        `json:"starts_at"  validate:"required"`
	EndsAt            time.Time        `json:"ends_at"    validate:"required"`
	PresenterName     *string          `json:"presenter_name"`
	TokenCost         int              `json:"token_cost"`
	HasCapacity       bool             `json:"has_capacity"`
	Capacity          int              `json:"capacity"`
	RemainingCapacity int              `json:"remaining_capacity"`
	Difficulty        *DifficultyLevel `json:"difficulty" validate:"oneof=no_prerequisites beginner intermediate advanced expert"`
	CreatedBy         uuid.UUID        `json:"created_by" validate:"required"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         *time.Time       `json:"deleted_at"`
}

type CreateActivitySpec struct {
	EditionID     uuid.UUID        `json:"edition_id"`
	Title         string           `json:"title"`
	Description   *string          `json:"description"`
	Location      string           `json:"location"`
	StartsAt      time.Time        `json:"starts_at"`
	EndsAt        time.Time        `json:"ends_at"`
	PresenterName *string          `json:"presenter_name"`
	TokenCost     int              `json:"token_cost"`
	HasCapacity   bool             `json:"has_capacity"`
	Capacity      int              `json:"capacity"`
	Difficulty    *DifficultyLevel `json:"difficulty"`
}

func NewActivity(creatorID uuid.UUID, spec CreateActivitySpec, edition *Edition) (*Activity, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("activity").SetMessage("error generating uuid").SetCause(err)
	}

	e := &Activity{
		ID:                id,
		EditionID:         spec.EditionID,
		Title:             spec.Title,
		Description:       spec.Description,
		Status:            ActivityStatusDraft,
		Location:          spec.Location,
		StartsAt:          spec.StartsAt,
		EndsAt:            spec.EndsAt,
		PresenterName:     spec.PresenterName,
		TokenCost:         spec.TokenCost,
		HasCapacity:       spec.HasCapacity,
		Capacity:          spec.Capacity,
		RemainingCapacity: spec.Capacity,
		Difficulty:        spec.Difficulty,
		CreatedBy:         creatorID,
	}

	if err := e.validate(edition); err != nil {
		return nil, err
	}

	return e, nil
}

func (a *Activity) validate(edition *Edition) error {
	now := time.Now()
	return validation.Run(
		validation.RequireUUID("activity", "edition_id", a.EditionID),
		validation.RequireString("activity", "title", a.Title),
		validation.RequireTime("activity", "starts_at", a.StartsAt),
		validation.RequireTime("activity", "ends_at", a.EndsAt),
		validation.Assert("activity", a.StartsAt.After(now), "start at must not be before now, legacy activities are not supported for now"),
		validation.Assert("activity", a.EndsAt.After(now), "ends at must not be before now, legacy activities are not supported for now"),
		validation.Assert("activity", a.StartsAt.Before(a.EndsAt), "starts must be before ends"),
		validation.Assert("activity", !a.StartsAt.Before(edition.StartsAt) && !a.StartsAt.After(edition.EndsAt), "activity start must be within edition duration"),
		validation.Assert("activity", !a.EndsAt.Before(edition.StartsAt) && !a.EndsAt.After(edition.EndsAt), "activity end must be within edition duration"),
		validation.Assert("activity", a.TokenCost >= 0, "invalid token cost amount"),
		validation.AssertIf("activity",
			func() bool { return a.HasCapacity },
			func() bool { return a.Capacity > 0 },
			"activity must have at least 1 capacity",
		),
		validation.AssertIf("activity",
			func() bool { return a.HasCapacity },
			func() bool { return a.Capacity == a.RemainingCapacity },
			"remaining capacity must be equal to capacity on creation",
		),
		validation.AssertIf("activity",
			func() bool { return a.HasCapacity },
			func() bool { return a.Capacity >= 0 },
			"invalid capacity amount",
		),
	)
}

func (e *Activity) AddScope(scopeID uuid.UUID) {
	e.ScopeID = scopeID
}

const AsynqActivityStart = "activity:start"
const AsynqActivityEnd = "activity:end"

type ActivityPayload struct {
	ActivityID uuid.UUID `json:"activity_id"`
}

func NewStartActivityTask(activityID uuid.UUID, startAt time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(ActivityPayload{ActivityID: activityID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(AsynqActivityStart, payload,
		asynq.TaskID(fmt.Sprintf("%s:%s", activityID, AsynqActivityStart)),
		asynq.ProcessAt(startAt),
		asynq.Unique(time.Hour),
		asynq.Retention(7*24*time.Hour),
	), nil
}

func NewEndActivityTask(activityID uuid.UUID, endAt time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(ActivityPayload{ActivityID: activityID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(AsynqActivityEnd, payload,
		asynq.TaskID(fmt.Sprintf("%s:%s", activityID, AsynqActivityEnd)),
		asynq.ProcessAt(endAt),
		asynq.Unique(time.Hour),
		asynq.Retention(7*24*time.Hour),
	), nil
}

type AttendanceStatus string

const (
	AttendanceStatusRegistered AttendanceStatus = "registered"
	AttendanceStatusWaitlisted AttendanceStatus = "waitlisted"
	AttendanceStatusPromoted   AttendanceStatus = "promoted"
	AttendanceStatusCheckedIn  AttendanceStatus = "checked_in"
	AttendanceStatusCheckedOut AttendanceStatus = "checked_out"
	AttendanceStatusCompleted  AttendanceStatus = "completed"
	AttendanceStatusPartial    AttendanceStatus = "partial"
	AttendanceStatusNoShow     AttendanceStatus = "no_show"
	AttendanceStatusCancelled  AttendanceStatus = "cancelled"
)

type AttendanceRecord struct {
	ID          uuid.UUID        `json:"id"`
	ActivityID  uuid.UUID        `json:"activity_id"`
	UserID      uuid.UUID        `json:"user_id"`
	Status      AttendanceStatus `json:"status"`
	CheckedInAt *time.Time       `json:"checked_in_at"`
	CancelledAt *time.Time       `json:"cancelled_at"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   *time.Time       `json:"deleted_at"`
}

type CreateAttendanceRecordSpec struct {
	ActivityID uuid.UUID        `json:"activity_id"`
	UserID     uuid.UUID        `json:"user_id"`
	Status     AttendanceStatus `json:"status"`
}

func NewAttendanceRecord(userID, activityID uuid.UUID) *AttendanceRecord {
	return &AttendanceRecord{
		ActivityID: activityID,
		UserID:     userID,
		Status:     AttendanceStatusRegistered,
	}
}

type CreateActivityRequest struct {
	Title         string           `json:"title" validate:"required,min=3"`
	Description   *string          `json:"description"`
	Location      string           `json:"location"`
	StartsAt      time.Time        `json:"starts_at" validate:"required"`
	EndsAt        time.Time        `json:"ends_at" validate:"required"`
	PresenterName *string          `json:"presenter_name"`
	TokenCost     int              `json:"token_cost" validate:"gte=0"`
	HasCapacity   bool             `json:"has_capacity"`
	Capacity      int              `json:"capacity" validate:"gte=0"`
	Difficulty    *DifficultyLevel `json:"difficulty"`
}

func (r CreateActivityRequest) ToSpec(editionID uuid.UUID) CreateActivitySpec {
	return CreateActivitySpec{
		EditionID:     editionID,
		Title:         r.Title,
		Description:   r.Description,
		Location:      r.Location,
		StartsAt:      r.StartsAt,
		EndsAt:        r.EndsAt,
		PresenterName: r.PresenterName,
		TokenCost:     r.TokenCost,
		HasCapacity:   r.HasCapacity,
		Capacity:      r.Capacity,
		Difficulty:    r.Difficulty,
	}
}
