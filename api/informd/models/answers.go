package models

import (
	"Informd/internal/database/sqlc"
	"encoding/json"
	"lib/xslices"
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

type SubmitAnswer struct {
	FieldID *uuid.UUID       `json:"field_id"`
	Answer  *json.RawMessage `json:"answer"`
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

func toSubmitAnswerInput(a SubmitAnswer) SubmitAnswerInput {
	return SubmitAnswerInput{
		FieldID: a.FieldID,
		Answer:  a.Answer,
	}
}

type SubmitRequest struct {
	Email   *string        `json:"email"`
	Answers []SubmitAnswer `json:"answers"`
}

func (r SubmitRequest) ToInput(formID uuid.UUID) SubmitInput {
	return SubmitInput{
		FormID:  formID,
		Email:   r.Email,
		Answers: xslices.MapSlice(r.Answers, toSubmitAnswerInput),
	}
}

type SubmitInput struct {
	FormID  uuid.UUID           `json:"form_id"`
	Email   *string             `json:"email"`
	Answers []SubmitAnswerInput `json:"answers"`
}
