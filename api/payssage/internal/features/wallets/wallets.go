// Package wallets is the wallets feature: its repo and handlers live here,
// its operations in service/, and Build wires the feature from here.
package wallets

import (
	"payssage/internal/features/wallets/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
