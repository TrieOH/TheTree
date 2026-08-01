package authn

import (
	"context"
	"lib/oauth"

	"github.com/MintzyG/fun"
)

func (o *Operations) OAuthConnect(_ context.Context, provider string) (string, error) {
	p, ok := oauth.Registry[provider]
	if !ok {
		return "", fun.ErrBadRequest("unsupported provider: " + provider)
	}
	return p.Config.AuthCodeURL(""), nil
}
