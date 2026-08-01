// Package responses implements the StrictServerInterface methods for the
// responses feature. Submission is public.
package responses

import (
	"context"
	"encoding/json"

	"Informd/internal/openapi"
	"Informd/internal/services"
	"Informd/models"
)

type Handlers struct {
	ops *services.Responses
}

func New(ops *services.Responses) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) SubmitResponse(ctx context.Context, req openapi.SubmitResponseRequestObject) (openapi.SubmitResponseResponseObject, error) {
	answers := make([]models.SubmitAnswerInput, 0, len(*req.Body.Answers))
	for _, a := range *req.Body.Answers {
		raw := jsonMarshal(a.Answer)
		answers = append(answers, models.SubmitAnswerInput{
			FieldID: &a.FieldId,
			Answer:  raw,
		})
	}
	err := h.ops.Submit(ctx, models.SubmitInput{
		FormID:  req.FormId,
		Email:   req.Body.Email,
		Answers: answers,
	})
	if err != nil {
		return nil, err
	}
	return openapi.SubmitResponse201Response{}, nil
}

func jsonMarshal(v *map[string]any) *json.RawMessage {
	if v == nil {
		return nil
	}
	b, _ := json.Marshal(v)
	raw := json.RawMessage(b)
	return &raw
}
