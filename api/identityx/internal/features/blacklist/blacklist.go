// Package blacklist is the blacklist feature: the layers live in the
// subpackages, and Build wires the feature from here.
package blacklist

import (
	"IdentityX/internal/features/blacklist/repo"
	"IdentityX/internal/sqlc"
)

type (
	Repo = repo.Repo
)

func NewRepo(q *sqlc.Queries) *Repo { return repo.New(q) }
