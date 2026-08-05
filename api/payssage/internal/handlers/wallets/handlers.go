package wallets

import (
	"payssage/internal/services"
)

type Handlers struct {
	ops *services.Wallets
}

func New(ops *services.Wallets) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
