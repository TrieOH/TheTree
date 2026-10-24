// Package action_tokens is the action tokens feature: the layers live in the
// subpackages, and Build wires the feature from here.
package action_tokens

import (
	"IdentityX/internal/features/action_tokens/repo"
	"IdentityX/internal/sqlc"
)

type (
	Repo = repo.Repo
)

func NewRepo(q *sqlc.Queries) *Repo { return repo.New(q) }
