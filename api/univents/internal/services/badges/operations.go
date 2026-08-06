package badges

import (
	"lib/email"
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	repo          ports.BadgeTemplateRepo
	emissions     ports.BadgeEmissionRepo
	registrations ports.RegistrationRepo
	editions      ports.EditionRepo
	events        ports.EventRepo
	email         *email.Client
	authz         *authz.Service
}

func NewOperations(
	repo ports.BadgeTemplateRepo,
	emissions ports.BadgeEmissionRepo,
	registrations ports.RegistrationRepo,
	editions ports.EditionRepo,
	events ports.EventRepo,
	email *email.Client,
	authz *authz.Service,
) *Operations {
	return &Operations{
		repo:          repo,
		emissions:     emissions,
		registrations: registrations,
		editions:      editions,
		events:        events,
		email:         email,
		authz:         authz,
	}
}
