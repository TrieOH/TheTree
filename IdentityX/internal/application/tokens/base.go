package tokens

import (
	"GoAuth/internal/adapters/persistence"
	"GoAuth/internal/application/tokens/issuer"
	"GoAuth/internal/application/tokens/verifier"
	"GoAuth/internal/ports/inbounds"
)

type TokenBundle struct {
	Issuer   inbounds.TokenIssuer
	Verifier inbounds.TokenVerifier
}

func NewBundle(repos *persistence.Repositories) TokenBundle {
	return TokenBundle{
		Issuer:   issuer.NewTokenIssuer(),
		Verifier: verifier.NewTokenVerifier(repos.Projects),
	}
}
