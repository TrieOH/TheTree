package fields

import (
	"encoding/json"

	"context"
	"time"

	"github.com/MintzyG/fun"

	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) EditSelectConfigNamespaced(ctx context.Context, req openapi.EditSelectConfigNamespacedRequestObject) (openapi.EditSelectConfigNamespacedResponseObject, error) {
	options, err := json.Marshal(req.Body.Options)
	if err != nil {
		return nil, fun.ErrBadRequest("invalid options")
	}
	config, err := h.ops.EditSelectConfigNamespaced(ctx, req.FormId, req.NamespaceId, models.FieldSelectConfig{
		FieldID:   req.FieldId,
		Behaviour: req.Body.Behaviour,
		ValueType: req.Body.ValueType,
		Options:   options,
	})
	if err != nil {
		return nil, err
	}
	return openapi.EditSelectConfigNamespaced200JSONResponse{
		Code: 200, Data: config, Timestamp: time.Now(), Module: module,
	}, nil
}
