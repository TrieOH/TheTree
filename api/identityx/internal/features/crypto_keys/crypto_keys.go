// Package crypto_keys is the crypto keys feature: the layers live in the
// subpackages, and Build wires the feature from here.
package crypto_keys

import (
	"IdentityX/internal/features/crypto_keys/repo"
	"IdentityX/internal/sqlc"
)

type (
	Repo = repo.Repo
)

func NewRepo(q *sqlc.Queries) *Repo { return repo.New(q) }
