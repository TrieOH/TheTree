package certifications

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.Certs
}

func New(ops *services.Certs) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
