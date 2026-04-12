package persistence

import (
	"GoAuth/internal/infrastructure"
	"GoAuth/internal/ports/outbounds"
)

type Repositories struct {
	Users          outbounds.UserRepository
	Sessions       outbounds.SessionRepository
	Projects       outbounds.ProjectRepository
	ProjectUsers   outbounds.ProjectUserRepository
	Keys           outbounds.KeysRepository
	TokenReuseList outbounds.TokenReuseListRepository
	ApiKey         outbounds.ApiKeyRepository
}

func NewRepositories(infra infrastructure.Infra) *Repositories {
	return &Repositories{
		Users:          NewUserRepo(infra.Queries, infra.Logger, infra.Tracer),
		Sessions:       NewSessionRepo(infra.Queries, infra.Logger, infra.Tracer),
		Projects:       NewProjectRepo(infra.Queries, infra.Logger, infra.Tracer),
		ProjectUsers:   NewProjectUserRepo(infra.Queries, infra.Logger, infra.Tracer),
		Keys:           NewKeyRepo(infra.Queries, infra.Logger, infra.Tracer),
		TokenReuseList: NewTokenReuseListRepo(infra.Queries, infra.Logger, infra.Tracer),
		ApiKey:         NewApiKeyRepo(infra.Queries, infra.Logger, infra.Tracer),
	}
}
