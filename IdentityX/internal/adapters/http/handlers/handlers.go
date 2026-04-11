package handlers

import (
	"GoAuth/internal/application"
	"GoAuth/internal/ports/outbounds"
)

type HandlerBundle struct {
	AuthHandler       *AuthHandler
	ProjectHandler    *ProjectHandler
	SessionHandler    *SessionHandler
	ScopeHandler      *ScopeHandler
	PermissionHandler *PermissionHandler
	RoleHandler       *RoleHandler
	ApiKeyHandler     *ApiKeyHandler
	SystemHandler     *SystemHandler
}

func New(app *application.Application, rdc outbounds.RedisCacheService) *HandlerBundle {
	return &HandlerBundle{
		AuthHandler:       NewAuthHandler(app.Auth, rdc),
		ProjectHandler:    NewProjectHandler(app.Project),
		SessionHandler:    NewSessionHandler(app.Session, rdc),
		ScopeHandler:      NewScopeHandler(app.Scope),
		PermissionHandler: NewPermissionHandler(app.Permission),
		RoleHandler:       NewRoleHandler(app.Role),
		ApiKeyHandler:     NewApiKeyHandler(app.ApiKey),
		SystemHandler:     NewSystemHandler(),
	}
}
