// Package oauth_providers owns everything about OAuth providers: the
// configured/enabled policy (CRUD, discovery), provider credential
// resolution, and the login flow (connect, callback, one-time login state,
// userinfo fetch, identity linking) at one interface. authn crosses it for
// the flow; handlers cross it for policy. Provider metadata and the HTTP
// client are injected at construction, so tests point both at a local
// server and never mutate a global registry.
package service

import (
	"IdentityX/internal/authz"
	tos "IdentityX/internal/features/tos/service"
	"IdentityX/internal/tokens"
	"IdentityX/ports"
	"lib/database"
	"lib/oauth"

	"resty.dev/v3"
)

type Deps struct {
	Providers   ports.ProjectOAuthProvidersRepo
	LoginStates ports.OAuthLoginStatesRepo
	Projects    ports.ProjectRepo
	External    ports.ExternalIdentitiesRepo
	Actors      ports.ActorRepo
	Authz       *authz.Service
	// Tokens owns the JWT lifecycle; the callback crosses it to mint the
	// session pair for a signed-in actor.
	Tokens *tokens.Manager
	// Tos owns the terms documents and the consent ledger; connect gates
	// on it and first-time registrations pin the acceptance row through it.
	Tos *tos.Operations
	// Client and meta are the HTTP seam: the token exchange and userinfo
	// calls go through the injected resty client, and provider metadata
	// (endpoints, userinfo URL) comes from the injected registry.
	Client *resty.Client
	Meta   map[string]oauth.Provider
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
