package events

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.Events
}

func New(ops *services.Events) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
