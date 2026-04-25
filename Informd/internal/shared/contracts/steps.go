package contracts

import (
	"github.com/google/uuid"
)

type Step struct {
	ID           uuid.UUID `json:"id"`
	FormID       uuid.UUID `json:"form_id" validate:"required"`
	Title        string    `json:"title" validate:"required"`
	Description  *string   `json:"description"`
	PositionHint int       `json:"position_hint" validate:"required"`
}

func NewStep(formID uuid.UUID, title string, description *string, positionHint int) (*Step, error) {
	s := &Step{
		FormID:       formID,
		Title:        title,
		Description:  description,
		PositionHint: positionHint,
	}
	if err := validate.Struct(s); err != nil {
		return nil, err
	}
	return s, nil
}
