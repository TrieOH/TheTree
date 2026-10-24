package service

import (
	"IdentityX/internal/authz"
	"IdentityX/internal/keys"
	"IdentityX/ports"
	"lib/database"
)

type Deps struct {
	// Keys owns the Key-lifecycle: project creation crosses its Ensure
	// seam so every project ships with signing and encryption keys.
	Keys     *keys.Manager
	Projects ports.ProjectRepo
	Actors   ports.ActorRepo
	Authz    *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
