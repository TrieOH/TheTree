package commands

import (
	"context"
	"lib/oauth"

	"github.com/MintzyG/fun"
)

func (c *Commands) OAuthConnect(ctx context.Context, provider string) (string, error) {
	p, ok := oauth.Registry[provider]
	if !ok {
		return "", fun.ErrBadRequest("unsupported provider: " + provider)
	}
	return p.Config.AuthCodeURL(""), nil
}
