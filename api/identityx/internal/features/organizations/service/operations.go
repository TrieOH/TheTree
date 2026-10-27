package service

import (
	"IdentityX/internal/authz"
	"IdentityX/internal/keys"
	"IdentityX/ports"
	"lib/database"
)

type Deps struct {
	Projects ports.ProjectRepo
	Actors   ports.ActorRepo
	Orgs     ports.OrganizationRepo
	// Keys owns the Key-lifecycle: org-created projects cross its Ensure
	// seam so they ship with keys instead of being token-broken until the
	// next boot.
	Keys  *keys.Manager
	Authz *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
