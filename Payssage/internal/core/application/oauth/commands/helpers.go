package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/errx"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (uc *CommandService) getProvider(name string) (domain.OAuthProvider, error) {
	p, ok := uc.providers[name]
	if !ok {
		return nil, errx.Invalid("provider").SetMessage(fmt.Sprintf("unsupported provider: %s", name))
	}
	return p, nil
}
