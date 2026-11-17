package fields

import (
	"context"
	"encoding/json"

	"Informd/errs"
	"Informd/internal/features/fields/service"
	"Informd/internal/openapi"
	"Informd/models"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func mustMarshal(v *map[string]any) *json.RawMessage {
	if v == nil {
		return nil
	}
	b, _ := json.Marshal(v)
	raw := json.RawMessage(b)
	return &raw
}

func mapSelectConfig(sc *openapi.CreateFieldSelectConfigRequest) *models.CreateFieldSelectConfigInput {
	if sc == nil {
		return nil
	}
	options, _ := json.Marshal(sc.Options)
	return &models.CreateFieldSelectConfigInput{
		Behaviour: sc.Behaviour,
		ValueType: sc.ValueType,
		Options:   options,
	}
}

func (h *Handlers) BulkEditFields(ctx context.Context, req openapi.BulkEditFieldsRequestObject) (openapi.BulkEditFieldsResponseObject, error) {
	payload := make([]models.UpdateStepFieldInput, 0, len(*req.Body))
	for _, f := range *req.Body {
		payload = append(payload, models.UpdateStepFieldInput{
			StepID:       req.StepId,
			ID:           f.Id,
			Key:          f.Key,
			Title:        f.Title,
			Description:  f.Description,
			PositionHint: f.PositionHint,
			Required:     f.Required,
			Type:         f.Type,
			Placeholder:  mustMarshal(f.Placeholder),
			DefaultValue: mustMarshal(f.DefaultValue),
			Config:       mustMarshal(f.Config),
			SelectConfig: mapSelectConfig(f.SelectConfig),
		})
	}
	err := h.ops.BulkEdit(ctx, req.FormId, payload)
	if err != nil {
		return nil, err
	}
	return openapi.BulkEditFields204Response{}, nil
}

func (h *Handlers) CreateField(ctx context.Context, req openapi.CreateFieldRequestObject) (openapi.CreateFieldResponseObject, error) {
	field, err := h.ops.Create(ctx, models.CreateStepFieldInput{
		FormID:       req.FormId,
		StepID:       req.StepId,
		Key:          req.Body.Key,
		Title:        req.Body.Title,
		Description:  req.Body.Description,
		PositionHint: req.Body.PositionHint,
		Required:     req.Body.Required,
		Type:         req.Body.Type,
		Placeholder:  mustMarshal(req.Body.Placeholder),
		DefaultValue: mustMarshal(req.Body.DefaultValue),
		Config:       mustMarshal(req.Body.Config),
		SelectConfig: mapSelectConfig(req.Body.SelectConfig),
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateField201JSONResponse(*field), nil
}

func (h *Handlers) DeleteField(ctx context.Context, req openapi.DeleteFieldRequestObject) (openapi.DeleteFieldResponseObject, error) {
	err := h.ops.Delete(ctx, req.FormId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteField204Response{}, nil
}

func (h *Handlers) EditSelectConfig(ctx context.Context, req openapi.EditSelectConfigRequestObject) (openapi.EditSelectConfigResponseObject, error) {
	options, err := json.Marshal(req.Body.Options)
	if err != nil {
		return nil, errs.ErrInvalidSelectOptions
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
	return openapi.EditSelectConfig200JSONResponse(*config), nil
}

func (h *Handlers) GetSelectConfig(ctx context.Context, req openapi.GetSelectConfigRequestObject) (openapi.GetSelectConfigResponseObject, error) {
	config, err := h.ops.GetSelectConfig(ctx, req.FormId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSelectConfig200JSONResponse(*config), nil
}

func (h *Handlers) ListFields(ctx context.Context, req openapi.ListFieldsRequestObject) (openapi.ListFieldsResponseObject, error) {
	fields, err := h.ops.List(ctx, req.FormId, req.StepId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFields200JSONResponse{Items: fields}, nil
}
