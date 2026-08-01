package orgs

import (
	"payssage/internal/services"
)

type Handlers struct {
	ops *services.Organizations
}

func New(ops *services.Organizations) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
