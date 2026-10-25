// Package platform_roles is the platform roles feature: the layers live in the
// subpackages, and Build wires the feature from here.
package platform_roles

import (
	"IdentityX/internal/features/platform_roles/repo"
	"IdentityX/internal/sqlc"
)

type (
	Repo = repo.Repo
)

func NewRepo(q *sqlc.Queries) *Repo { return repo.New(q) }
