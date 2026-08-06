package models

import (
	"context"

	"github.com/MintzyG/fun"
)

// RequireClientOnly rejects unauthenticated requests and requests from
// project-scoped actors. It requires a valid identity in the context whose
// subject carries a nil ProjectID — an IdentityX platform-level client
// (human, service, or machine) rather than a project-scoped one.
//
// Call it at the top of handlers that should only be reachable by
// platform-level IdentityX clients. It replaces the old ClientOnly route
// middleware: the scope policy now lives at the handler that enforces it,
// not in the auth chain.
func RequireClientOnly(ctx context.Context) error {
	ident, err := RequireIdentity(ctx)
	if err != nil {
		return err
	}
	if ident.Sub.ProjectID != nil {
		return fun.ErrUnauthorized("platform-level authentication required")
	}
	return nil
}

// RequireProjectClientOnly rejects requests that are not both authenticated
// AND scoped to a specific project. It requires a valid identity whose
// subject carries a non-nil ProjectID — an actor acting within a project
// context (e.g. workspace API keys, project service accounts).
//
// Not used yet; defined for the project-scoped operations that will need
// it. Mirror of RequireClientOnly.
func RequireProjectClientOnly(ctx context.Context) error {
	ident, err := RequireIdentity(ctx)
	if err != nil {
		return err
	}
	if ident.Sub.ProjectID == nil {
		return fun.ErrUnauthorized("project-scoped authentication required")
	}
	return nil
}
