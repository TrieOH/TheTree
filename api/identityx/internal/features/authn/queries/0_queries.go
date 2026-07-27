package queries

import (
	"IdentityX/ports"
)

type Queries struct {
	projects   ports.ProjectRepo
	cryptoKeys ports.CryptoKeysRepo
}

func NewQueries(
	projects ports.ProjectRepo,
	cryptoKeys ports.CryptoKeysRepo,
) *Queries {
	return &Queries{
		projects:   projects,
		cryptoKeys: cryptoKeys,
	}
}
