package authn

import (
	"context"
	"lib/oauth"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) OAuthConnect(ctx context.Context, provider string, projectID *uuid.UUID) (string, error) {
	ctx, span := telemetry.StartSpan(ctx, "OAuthConnect")
	defer span.End()

	p, ok := oauth.Registry[provider]
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
	}

	// Provider policy lives in the oauth_providers module: this is the
	// flow's only contact with "is this provider configured/enabled".
	resolved, err := o.oauthProviders.ResolveLoginProvider(ctx, provider, projectID)
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
