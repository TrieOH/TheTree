// Package fields implements the StrictServerInterface methods for the
// fields feature, including the duplicated namespaced routes.
package fields

import (
	"context"
	"encoding/json"
	"time"

	"Informd/internal/openapi"
	"Informd/internal/services"
	"Informd/models"

	"github.com/MintzyG/fun"
)

const module = "Informd"

type Handlers struct {
	ops *services.Fields
}

func New(ops *services.Fields) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListFields(ctx context.Context, req openapi.ListFieldsRequestObject) (openapi.ListFieldsResponseObject, error) {
	fields, err := h.ops.List(ctx, req.FormId, req.StepId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFields200JSONResponse{
		Code: 200, Data: &fields, Timestamp: time.Now(), Module: module,
	}, nil
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
	return openapi.CreateField201JSONResponse{
		Code: 201, Data: field, Timestamp: time.Now(), Module: module,
	}, nil
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
	return openapi.BulkEditFields200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetSelectConfig(ctx context.Context, req openapi.GetSelectConfigRequestObject) (openapi.GetSelectConfigResponseObject, error) {
	config, err := h.ops.GetSelectConfig(ctx, req.FormId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSelectConfig200JSONResponse{
		Code: 200, Data: config, Timestamp: time.Now(), Module: module,
	}, nil
}

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

func (h *Handlers) DeleteField(ctx context.Context, req openapi.DeleteFieldRequestObject) (openapi.DeleteFieldResponseObject, error) {
	err := h.ops.Delete(ctx, req.FormId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteField200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListFieldsNamespaced(ctx context.Context, req openapi.ListFieldsNamespacedRequestObject) (openapi.ListFieldsNamespacedResponseObject, error) {
	fields, err := h.ops.ListNamespaced(ctx, req.FormId, req.NamespaceId, req.StepId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFieldsNamespaced200JSONResponse{
		Code: 200, Data: &fields, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CreateFieldNamespaced(ctx context.Context, req openapi.CreateFieldNamespacedRequestObject) (openapi.CreateFieldNamespacedResponseObject, error) {
	field, err := h.ops.CreateNamespaced(ctx, models.CreateNamespacedStepFieldInput{
		NamespaceID:  req.NamespaceId,
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
	return openapi.CreateFieldNamespaced201JSONResponse{
		Code: 201, Data: field, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) BulkEditFieldsNamespaced(ctx context.Context, req openapi.BulkEditFieldsNamespacedRequestObject) (openapi.BulkEditFieldsNamespacedResponseObject, error) {
	payload := make([]models.UpdateNamespacedStepFieldInput, 0, len(*req.Body))
	for _, f := range *req.Body {
		payload = append(payload, models.UpdateNamespacedStepFieldInput{
			NamespaceID:  req.NamespaceId,
			FormID:       req.FormId,
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
	err := h.ops.BulkEditNamespaced(ctx, req.FormId, req.NamespaceId, payload)
	if err != nil {
		return nil, err
	}
	return openapi.BulkEditFieldsNamespaced200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetSelectConfigNamespaced(ctx context.Context, req openapi.GetSelectConfigNamespacedRequestObject) (openapi.GetSelectConfigNamespacedResponseObject, error) {
	config, err := h.ops.GetSelectConfigNamespaced(ctx, req.FormId, req.NamespaceId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSelectConfigNamespaced200JSONResponse{
		Code: 200, Data: config, Timestamp: time.Now(), Module: module,
	}, nil
}

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

func (h *Handlers) DeleteFieldNamespaced(ctx context.Context, req openapi.DeleteFieldNamespacedRequestObject) (openapi.DeleteFieldNamespacedResponseObject, error) {
	err := h.ops.DeleteNamespaced(ctx, req.NamespaceId, req.FormId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteFieldNamespaced200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}

func mustMarshal(v *map[string]any) *json.RawMessage {
	if v == nil {
		return nil
	}
	b, _ := json.Marshal(v)
	raw := json.RawMessage(b)
	return &raw
}

func mapSelectConfig(sc *openapi.CreateFieldSelectConfigRequest) *models.CreateFieldSelectConfigRequest {
	if sc == nil {
		return nil
	}
	options, _ := json.Marshal(sc.Options)
	return &models.CreateFieldSelectConfigRequest{
		Behaviour: sc.Behaviour,
		ValueType: sc.ValueType,
		Options:   options,
	}
}
