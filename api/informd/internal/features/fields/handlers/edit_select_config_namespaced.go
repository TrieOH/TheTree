package handlers

import (
	"encoding/json"

	"context"

	"Informd/errs"
	"Informd/internal/openapi"
	"Informd/models"
)

func (h *Handlers) EditSelectConfigNamespaced(ctx context.Context, req openapi.EditSelectConfigNamespacedRequestObject) (openapi.EditSelectConfigNamespacedResponseObject, error) {
	options, err := json.Marshal(req.Body.Options)
	if err != nil {
		return nil, errs.ErrInvalidSelectOptions
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
	return openapi.EditSelectConfigNamespaced200JSONResponse(*config), nil
}
