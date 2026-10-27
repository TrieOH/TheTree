// Package registrations is the registrations feature: the layers live in the
// subpackages, and Build wires the feature from here.
package registrations

import (
	"univents/internal/features/registrations/repo"
	"univents/internal/sqlc"
)

type (
	Repo = repo.Repo
)

func NewRepo(q *sqlc.Queries) *Repo { return repo.New(q) }
