// Package events is the events feature: its repo and handlers live here,
// its operations in service/, and Build wires the feature from here.
package events

import (
	"univents/internal/features/events/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
