package responses

import (
	"context"

	"Informd/internal/openapi"
	"Informd/models"
)

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
