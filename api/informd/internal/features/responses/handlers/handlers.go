package handlers

import (
	"encoding/json"

	"Informd/internal/features/responses/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func jsonMarshal(v *map[string]any) *json.RawMessage {
	if v == nil {
		return nil
	}
	b, _ := json.Marshal(v)
	raw := json.RawMessage(b)
	return &raw
}
