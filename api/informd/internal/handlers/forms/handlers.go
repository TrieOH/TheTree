package forms

import (
	"Informd/internal/services"
)

type Handlers struct {
	ops *services.Forms
}

func New(ops *services.Forms) *Handlers { return &Handlers{ops: ops} }

const module = "Informd"
