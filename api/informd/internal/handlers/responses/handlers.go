package responses

import (
	"encoding/json"

	"Informd/internal/services"
)

type Handlers struct {
	ops *services.Responses
}

func New(ops *services.Responses) *Handlers { return &Handlers{ops: ops} }

func jsonMarshal(v *map[string]any) *json.RawMessage {
	if v == nil {
		return nil
	}
	b, _ := json.Marshal(v)
	raw := json.RawMessage(b)
	return &raw
}
