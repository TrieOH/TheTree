package signatures

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.Signatures
}

func New(ops *services.Signatures) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
