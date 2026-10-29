// Package ticket_types is the ticket types feature: its repo and handlers
// live here, its operations in service/, and Build wires the feature from
// here.
package ticket_types

import (
	"univents/internal/features/ticket_types/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
