package service

import (
	"strings"

	"lib/crypto"
	"lib/oauth"
)

// testVault seals and opens test keys with a fixed master key.
var testVault = func() *crypto.Vault {
	v, err := crypto.NewVault(strings.Repeat("ab", 32))
	if err != nil {
		panic(err)
	}
	return v
}()

// platformGoogle configures IdentityX's own google login.
var platformGoogle = map[string]oauth.Credentials{
	"google": {ClientID: "platform-id", ClientSecret: "platform-secret"},
}
