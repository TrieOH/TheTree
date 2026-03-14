package handlers

import (
	"GoAuth/internal/application"
	"GoAuth/internal/ports/outbounds"
)

type HandlerBundle struct {
	AuthHandler          *AuthHandler
	ProjectHandler       *ProjectHandler
	SessionHandler       *SessionHandler
	SchemaHandler        *SchemaHandler
	SchemaVersionHandler *SchemaVersionHandler
	SchemaFieldsHandler  *SchemaFieldsHandler
	ScopeHandler         *ScopeHandler
	PermissionHandler    *PermissionHandler
	RoleHandler          *RoleHandler
	ApiKeyHandler        *ApiKeyHandler
	SubContextHandler    *SubContextHandler
	SystemHandler        *SystemHandler
}

func New(app *application.Application, rdc outbounds.RedisCacheService) *HandlerBundle {
	return &HandlerBundle{
		AuthHandler:          NewAuthHandler(app.Auth, app.Schema),
		ProjectHandler:       NewProjectHandler(app.Project),
		SessionHandler:       NewSessionHandler(app.Session, rdc),
		SchemaHandler:        NewSchemaHandler(app.Schema),
		SchemaVersionHandler: NewSchemaVersionHandler(app.SchemaVersions),
		SchemaFieldsHandler:  NewSchemaFieldsHandler(app.SchemaFields),
		ScopeHandler:         NewScopeHandler(app.Scope),
		PermissionHandler:    NewPermissionHandler(app.Permission),
		RoleHandler:          NewRoleHandler(app.Role),
		ApiKeyHandler:        NewApiKeyHandler(app.ApiKey),
		SubContextHandler:    NewSubContextHandler(app.SubContext),
		SystemHandler:        NewSystemHandler(),
	}
}
