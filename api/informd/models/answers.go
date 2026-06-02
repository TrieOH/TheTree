package models

import (
	"Informd/internal/database/sqlc"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         uuid.UUID        `json:"id"`
	ResponseID uuid.UUID        `json:"response_id"`
	FieldID    *uuid.UUID       `json:"field_id"`
	Answer     *json.RawMessage `json:"answer"`
	AnsweredAt time.Time        `json:"answered_at"`
	UpdatedAt  *time.Time       `json:"updated_at"`
}

func ToBatchUpsertAnswersParams(a Answer) sqlc.BatchUpsertAnswersParams {
	return sqlc.BatchUpsertAnswersParams{
		ResponseID: a.ResponseID,
		FieldID:    a.FieldID,
		Answer:     a.Answer,
	}
}

type SubmitRequest struct {
	Email   *string  `json:"email"`
	Answers []Answer `json:"answers"`
}

func (r SubmitRequest) ToInput(formID uuid.UUID) SubmitInput {
	return SubmitInput{
		FormID:  formID,
		Email:   r.Email,
		Answers: r.Answers,
	}
}

type SubmitInput struct {
	FormID  uuid.UUID `json:"form_id"`
	Email   *string   `json:"email"`
	Answers []Answer  `json:"answers"`
}
