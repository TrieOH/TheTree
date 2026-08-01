package models

import (
	"Informd/internal/sqlc"
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

type SubmitAnswerInput struct {
	FieldID    *uuid.UUID       `json:"field_id"`
	Answer     *json.RawMessage `json:"answer"`
	ResponseID uuid.UUID        `json:"response_id"`
}

func SubmitAnswerInputToAnswer(input SubmitAnswerInput) Answer {
	return Answer{
		FieldID:    input.FieldID,
		Answer:     input.Answer,
		ResponseID: input.ResponseID,
	}
}

type SubmitInput struct {
	FormID  uuid.UUID           `json:"form_id"`
	Email   *string             `json:"email"`
	Answers []SubmitAnswerInput `json:"answers"`
}
