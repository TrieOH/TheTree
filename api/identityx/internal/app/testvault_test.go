package app

import (
	"strings"

	"lib/crypto"
)

// testVault seals and opens test keys with a fixed master key.
var testVault = func() *crypto.Vault {
	v, err := crypto.NewVault(strings.Repeat("ab", 32))
	if err != nil {
		panic(err)
	}
	return v
}()
