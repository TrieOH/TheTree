package models

import (
	"Informd/internal/database/sqlc"

	"github.com/google/uuid"
)

type Step struct {
	ID           uuid.UUID `json:"id"`
	FormID       uuid.UUID `json:"form_id" validate:"required"`
	Title        string    `json:"title" validate:"required"`
	Description  *string   `json:"description"`
	PositionHint int       `json:"position_hint" validate:"required,gte=1"`
}

func NewStep(formID uuid.UUID, title string, description *string, positionHint int) (*Step, error) {
	s := &Step{
		FormID:       formID,
		Title:        title,
		Description:  description,
		PositionHint: positionHint,
	}
	return s, validate.Struct(s)
}

type CreateStepRequest struct {
	Title        string  `json:"title" validate:"required"`
	Description  *string `json:"description"`
	PositionHint int     `json:"position_hint" validate:"required"`
}

func (r CreateStepRequest) ToFormInput(formID uuid.UUID) CreateFormStepInput {
	return CreateFormStepInput{
		FormID:       formID,
		Title:        r.Title,
		Description:  r.Description,
		PositionHint: r.PositionHint,
	}
}

func (r CreateStepRequest) ToNamespacedFormInput(namespaceID, formID uuid.UUID) CreateNamespacedFormStepInput {
	return CreateNamespacedFormStepInput{
		NamespaceID:  namespaceID,
		FormID:       formID,
		Title:        r.Title,
		Description:  r.Description,
		PositionHint: r.PositionHint,
	}
}

type CreateFormStepInput struct {
	FormID       uuid.UUID `json:"form_id"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	PositionHint int       `json:"position_hint"`
}

type CreateNamespacedFormStepInput struct {
	NamespaceID  uuid.UUID `json:"namespace_id"`
	FormID       uuid.UUID `json:"form_id"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	PositionHint int       `json:"position_hint"`
}

func ToBulkEditStepsParams(s Step) sqlc.BulkEditStepsParams {
	return sqlc.BulkEditStepsParams{
		ID:           s.ID,
		FormID:       s.FormID,
		Title:        s.Title,
		Description:  s.Description,
		PositionHint: s.PositionHint,
	}
}

type UpdateStepRequest struct {
	ID           uuid.UUID `json:"id" validate:"required"`
	Title        string    `json:"title" validate:"required"`
	Description  *string   `json:"description"`
	PositionHint int       `json:"position_hint" validate:"required,gte=1"`
}

func (r UpdateStepRequest) ToFormInput(formID uuid.UUID) UpdateFormStepInput {
	return UpdateFormStepInput{
		FormID:       formID,
		ID:           r.ID,
		Title:        r.Title,
		Description:  r.Description,
		PositionHint: r.PositionHint,
	}
}

func (r UpdateStepRequest) ToNamespacedFormInput(namespaceID, formID uuid.UUID) UpdateNamespacedFormStepInput {
	return UpdateNamespacedFormStepInput{
		NamespaceID:  namespaceID,
		FormID:       formID,
		ID:           r.ID,
		Title:        r.Title,
		Description:  r.Description,
		PositionHint: r.PositionHint,
	}
}

type UpdateFormStepInput struct {
	FormID       uuid.UUID
	ID           uuid.UUID
	Title        string
	Description  *string
	PositionHint int
}

type UpdateNamespacedFormStepInput struct {
	NamespaceID  uuid.UUID
	FormID       uuid.UUID
	ID           uuid.UUID
	Title        string
	Description  *string
	PositionHint int
}

func UpdateFormStepInputToStep(i UpdateFormStepInput) Step {
	return Step{
		ID:           i.ID,
		FormID:       i.FormID,
		Title:        i.Title,
		Description:  i.Description,
		PositionHint: i.PositionHint,
	}
}

func UpdateNamespacedFormStepInputToStep(i UpdateNamespacedFormStepInput) Step {
	return Step{
		ID:           i.ID,
		FormID:       i.FormID,
		Title:        i.Title,
		Description:  i.Description,
		PositionHint: i.PositionHint,
	}
}
