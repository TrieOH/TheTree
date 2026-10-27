// Package ws_tokens is the ws tokens feature: the layers live in the
// subpackages, and Build wires the feature from here.
package ws_tokens

import (
	"univents/internal/features/ws_tokens/repo"
	"univents/internal/sqlc"
)

type (
	Repo = repo.Repo
)

func NewRepo(q *sqlc.Queries) *Repo { return repo.New(q) }
