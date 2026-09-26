package oauth_providers

import (
	"context"

	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// Connect starts an OAuth login: it resolves the scope's provider
// credentials, records a one-time login state, and returns the provider's
// authorization URL for the caller to redirect to. Provider policy lives
// in this module, so this is the flow's only contact with "is this
// provider configured/enabled".
//
// acceptedTos is the clickwrap consent flag: a project with terms only
// starts flows whose caller agreed to them, so the consent screen precedes
// the provider redirect (the callback itself has no UI to ask in).
func (o *Operations) Connect(ctx context.Context, provider string, projectID *uuid.UUID, acceptedTos bool) (string, error) {
	ctx, span := telemetry.StartSpan(ctx, "Connect")
	defer span.End()

	p, ok := o.meta[provider]
	if !ok {
		return "", fun.ErrBadRequest("unsupported provider: " + provider)
	}

	if projectID != nil {
		// The project must exist before its provider row can resolve; a
		// missing project surfaces as not-found, not as a misconfigured row.
		_, err := o.projects.GetByID(ctx, *projectID)
		if err != nil {
			return "", err
		}
		_, hasTos, err := o.tos.CurrentTos(ctx, *projectID)
		if err != nil {
			return "", err
		}
		if hasTos && !acceptedTos {
			return "", fun.ErrValidation("you must accept the terms of service to register")
		}
	}

	resolved, err := o.resolveLoginProvider(ctx, provider, projectID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return "", fun.ErrBadRequest("provider not configured for this project: " + provider)
		}
		return "", err
	}

	state, err := o.createLoginState(ctx, provider, projectID)
	if err != nil {
		return "", err
	}

	return p.Config(resolved.Creds).AuthCodeURL(state), nil
}
