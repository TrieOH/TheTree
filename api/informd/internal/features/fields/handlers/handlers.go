package handlers

import (
	"encoding/json"

	"Informd/internal/features/fields/service"
	"Informd/internal/openapi"
	"Informd/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

const module = "Informd"

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
