// Package responders is the responders feature: the layers live in the
// subpackages, and Build wires the feature from here.
package responders

import (
	"Informd/internal/features/responders/repo"
	"Informd/internal/sqlc"
)

type (
	Repo = repo.Repo
)

func NewRepo(q *sqlc.Queries) *Repo { return repo.New(q) }
