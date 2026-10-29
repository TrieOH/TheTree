// Package tos is the tos feature: its repo and handlers live here, its
// operations in service/, and Build wires the feature from here.
//
// The repo persists the terms-of-service documents and their consent
// ledger; the handlers serve the versioned document (get/post/put) and the
// consent ledger (accept/list).
package tos

import (
	"IdentityX/internal/features/tos/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
