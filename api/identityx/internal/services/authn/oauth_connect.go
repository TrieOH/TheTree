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

	creds, err := o.resolveCredentials(ctx, provider, projectID)
	if err != nil {
		return "", err
	}

	state, err := o.createLoginState(ctx, provider, projectID)
	if err != nil {
		return "", err
	}

	return p.Config(creds.creds).AuthCodeURL(state), nil
}
