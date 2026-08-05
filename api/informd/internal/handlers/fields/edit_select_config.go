package fields

import (
	"encoding/json"

	"context"
	"time"

	"github.com/MintzyG/fun"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) EditSelectConfig(ctx context.Context, req openapi.EditSelectConfigRequestObject) (openapi.EditSelectConfigResponseObject, error) {
	options, err := json.Marshal(req.Body.Options)
	if err != nil {
		return nil, fun.ErrBadRequest("invalid options")
	}
	config, err := h.ops.EditSelectConfig(ctx, req.FormId, models.FieldSelectConfig{
		FieldID:   req.FieldId,
		Behaviour: req.Body.Behaviour,
		ValueType: req.Body.ValueType,
		Options:   options,
	})
	if err != nil {
		return nil, err
	}
	return openapi.EditSelectConfig200JSONResponse{
		Code: 200, Data: config, Timestamp: time.Now(), Module: module,
	}, nil
}
