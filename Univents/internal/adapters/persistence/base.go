package persistence

import (
	"univents/internal/infrastructure"
)

type Repositories struct {
}

func NewRepositories(infra infrastructure.Infra) *Repositories {
	return &Repositories{}
}
