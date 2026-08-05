package badges

import (
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	repo     ports.BadgeTemplateRepo
	editions ports.EditionRepo
	authz    *authz.Service
}

func NewOperations(
	repo ports.BadgeTemplateRepo,
	editions ports.EditionRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		repo:     repo,
		editions: editions,
		authz:    authz,
	}
}
