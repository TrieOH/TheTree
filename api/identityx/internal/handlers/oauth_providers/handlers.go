package oauth_providers

import (
	"IdentityX/internal/services"
)

type Handlers struct {
	ops *services.OAuthProviders
}

func New(ops *services.OAuthProviders) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"
