package service

import (
	tos "IdentityX/internal/features/tos/service"
	"IdentityX/internal/tokens"
	"IdentityX/ports"
	"lib/database"
)

// Operations owns the session and account lifecycle: password login,
// registration, refresh, logout, setup, JWKS, and the verify/reset email
// flows. OAuth login lives in the oauth_providers module — authn crosses
// its Connect/Callback interface and touches no OAuth state itself.
type Deps struct {
	Actors        ports.ActorRepo
	Projects      ports.ProjectRepo
	PlatformRoles ports.PlatformRolesRepo
	// Tokens owns the token lifecycle (verify/mint/rotate/revoke); login,
	// refresh, logout cross it instead of touching keys, blacklist, or
	// token claims directly.
	Tokens *tokens.Manager
	// ActionTokens owns the single-use action-token lifecycle; verify and
	// reset links are redeemed through it instead of touching the HMAC
	// secret or the anti-replay repo directly.
	ActionTokens *tokens.ActionTokenManager
	EmailSender  ports.EmailSender
	// Tos owns the terms documents and the consent ledger; registration
	// gates on it and pins the acceptance row through it.
	Tos *tos.Operations
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
