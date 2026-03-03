package domain

import (
	"time"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
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
	ScopeID           uuid.UUID        `json:"scope_id"`
	EditionID         uuid.UUID        `json:"edition_id"`
	Title             string           `json:"title"`
	Description       *string          `json:"description"`
	Status            ActivityStatus   `json:"status"`
	Location          string           `json:"location"`
	StartsAt          time.Time        `json:"starts_at"`
	EndsAt            time.Time        `json:"ends_at"`
	PresenterName     *string          `json:"presenter_name"`
	TokenCost         int              `json:"token_cost"`
	HasCapacity       bool             `json:"has_capacity"`
	Capacity          int              `json:"capacity"`
	RemainingCapacity int              `json:"remaining_capacity"`
	Difficulty        *DifficultyLevel `json:"difficulty"`
	CreatedBy         uuid.UUID        `json:"created_by"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         *time.Time       `json:"deleted_at"`
}

type CreateActivitySpec struct {
	EditionScopeID uuid.UUID        `json:"edition_scope_id"`
	EditionID      uuid.UUID        `json:"edition_id"`
	Title          string           `json:"title"`
	Description    *string          `json:"description"`
	Location       string           `json:"location"`
	StartsAt       time.Time        `json:"starts_at"`
	EndsAt         time.Time        `json:"ends_at"`
	PresenterName  *string          `json:"presenter_name"`
	TokenCost      int              `json:"token_cost"`
	HasCapacity    bool             `json:"has_capacity"`
	Capacity       int              `json:"capacity"`
	Difficulty     *DifficultyLevel `json:"difficulty"`
}

func NewActivity(creatorID uuid.UUID, spec CreateActivitySpec, edition *Edition) (*Activity, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fail.New(errx.SYSUUIDV7GenerationError).WithArgs("NewActivity")
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

// FIXME make me unique errors
func (a *Activity) validate(edition *Edition) error {
	if a.EditionID == uuid.Nil {
		return fail.New(errx.ActivityValidationFailed).WithArgs(uuid.Nil.String())
	}

	if a.Title == "" {
		return fail.New(errx.ActivityValidationFailed).Trace("activity title is required")
	}

	if a.StartsAt.IsZero() || a.EndsAt.IsZero() {
		return fail.New(errx.EditionValidationFailed).Trace("start and end times are required")
	}

	now := time.Now()
	if a.StartsAt.Before(now) {
		return fail.New(errx.EditionValidationFailed).Trace("start at must not be before now, legacy activities are not supported for now")
	}
	if a.EndsAt.Before(now) {
		return fail.New(errx.EditionValidationFailed).Trace("start at must not be before now, legacy activities are not supported for now")
	}

	if a.StartsAt.Before(edition.StartsAt) || a.StartsAt.After(edition.EndsAt) {
		return fail.New(errx.EditionValidationFailed).Trace("activity start must be within edition duration")
	}

	if a.EndsAt.Before(edition.StartsAt) || a.EndsAt.After(edition.EndsAt) {
		return fail.New(errx.EditionValidationFailed).Trace("activity end must be within edition duration")
	}

	if !a.StartsAt.Before(a.EndsAt) {
		return fail.New(errx.EditionValidationFailed).Trace("starts_must_be_before_ends")
	}

	if a.HasCapacity {
		if a.Capacity == 0 {
			return fail.New(errx.EditionValidationFailed).Trace("activity must have at least 1 capacity")
		}
		if a.Capacity != a.RemainingCapacity {
			return fail.New(errx.EditionValidationFailed).Trace("remaining capacity must be equal to capacity on creation")
		}
		if a.Capacity < 0 {
			return fail.New(errx.EditionValidationFailed).Trace("invalid capacity amount")
		}
	}

	if a.TokenCost < 0 {
		return fail.New(errx.EditionValidationFailed).Trace("invalid token cost amount")
	}

	return nil
}

func (e *Activity) AddScope(scopeID uuid.UUID) {
	e.ScopeID = scopeID
}
