package service

import (
	"Informd/internal/authz"

	"Informd/ports"
)

type Deps struct {
	Forms      ports.FormsRepo
	Steps      ports.StepRepo
	Namespaces ports.NamespaceRepo
	Authz      *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
