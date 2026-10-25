// Package answers is the answers feature: the layers live in the
// subpackages, and Build wires the feature from here.
package answers

import (
	"Informd/internal/features/answers/repo"
	"Informd/internal/sqlc"
)

type (
	Repo = repo.Repo
)

func NewRepo(q *sqlc.Queries) *Repo { return repo.New(q) }
