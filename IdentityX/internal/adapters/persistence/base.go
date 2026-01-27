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
	Schemas        outbounds.SchemaRepository
	SchemaVersions outbounds.SchemaVersionRepository
	SchemaFields   outbounds.SchemaFieldsRepository
	Scopes         outbounds.ScopeRepository
	Permissions    outbounds.PermissionRepository
	Keys           outbounds.KeysRepository
}

func NewRepositories(infra infrastructure.Infra) *Repositories {
	return &Repositories{
		Users:          NewUserRepo(infra.Queries, infra.Logger, infra.Tracer),
		Sessions:       NewSessionRepo(infra.Queries, infra.Logger, infra.Tracer),
		Projects:       NewProjectRepo(infra.Queries, infra.Logger, infra.Tracer),
		ProjectUsers:   NewProjectUserRepo(infra.Queries, infra.Logger, infra.Tracer),
		Schemas:        NewSchemaRepo(infra.Queries, infra.Logger, infra.Tracer),
		SchemaVersions: NewSchemaVersionRepo(infra.Queries, infra.Logger, infra.Tracer),
		SchemaFields:   NewFieldsRepo(infra.Queries, infra.Logger, infra.Tracer),
		Scopes:         NewScopeRepo(infra.Queries, infra.Logger, infra.Tracer),
		Permissions:    NewPermissionRepo(infra.Queries, infra.Logger, infra.Tracer),
		Keys:           NewKeyRepo(infra.Queries, infra.Logger, infra.Tracer),
	}
}
