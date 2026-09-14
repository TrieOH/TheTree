package authn

import (
	"IdentityX/internal/services/tos"
	"IdentityX/internal/tokens"
	"IdentityX/ports"
	"lib/errx"
)

// Operations owns the session and account lifecycle: password login,
// registration, refresh, logout, setup, JWKS, and the verify/reset email
// flows. OAuth login lives in the oauth_providers module — authn crosses
// its Connect/Callback interface and touches no OAuth state itself.
type Operations struct {
	actors        ports.ActorRepo
	projects      ports.ProjectRepo
	platformRoles ports.PlatformRolesRepo
	// tokens owns the token lifecycle (verify/mint/rotate/revoke); login,
	// refresh, logout cross it instead of touching keys, blacklist, or
	// token claims directly.
	tokens *tokens.Manager
	// actionTokens owns the single-use action-token lifecycle; verify and
	// reset links are redeemed through it instead of touching the HMAC
	// secret or the anti-replay repo directly.
	actionTokens *tokens.ActionTokenManager
	emailSender  ports.EmailSender
	// tos owns the terms documents and the consent ledger; registration
	// gates on it and pins the acceptance row through it.
	tos *tos.Operations
}

func NewOperations(
	actors ports.ActorRepo,
	projects ports.ProjectRepo,
	platformRoles ports.PlatformRolesRepo,
	tokensMgr *tokens.Manager,
	actionTokens *tokens.ActionTokenManager,
	emailSender ports.EmailSender,
	tosOps *tos.Operations,
) *Operations {
	return errx.MustProvide(&Operations{
		actors:        actors,
		projects:      projects,
		platformRoles: platformRoles,
		tokens:        tokensMgr,
		actionTokens:  actionTokens,
		emailSender:   emailSender,
		tos:           tosOps,
	})
}
