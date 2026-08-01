package badges

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.Badges
}

func New(ops *services.Badges) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
