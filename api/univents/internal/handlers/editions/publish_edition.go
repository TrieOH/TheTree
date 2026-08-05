package editions

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) PublishEdition(ctx context.Context, req openapi.PublishEditionRequestObject) (openapi.PublishEditionResponseObject, error) {
	err := h.ops.Publish(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.PublishEdition204Response{}, nil
}
