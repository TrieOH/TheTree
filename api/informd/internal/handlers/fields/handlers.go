package fields

import (
	"encoding/json"

	"Informd/internal/openapi"
	"Informd/internal/services"
	"Informd/models"
)

type Handlers struct {
	ops *services.Fields
}

func New(ops *services.Fields) *Handlers { return &Handlers{ops: ops} }

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
