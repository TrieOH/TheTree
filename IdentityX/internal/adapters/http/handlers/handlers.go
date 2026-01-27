package handlers

import (
	"GoAuth/internal/application"
)

type HandlerBundle struct {
	AuthHandler          *AuthHandler
	ProjectHandler       *ProjectHandler
	SessionHandler       *SessionHandler
	SchemaHandler        *SchemaHandler
	SchemaVersionHandler *SchemaVersionHandler
	SchemaFieldsHandler  *SchemaFieldsHandler
	ScopeHandler         *ScopeHandler
}

func New(app *application.Application) *HandlerBundle {
	return &HandlerBundle{
		AuthHandler:          NewAuthHandler(app.Auth),
		ProjectHandler:       NewProjectHandler(app.Project),
		SessionHandler:       NewSessionHandler(app.Session),
		SchemaHandler:        NewSchemaHandler(app.Schema),
		SchemaVersionHandler: NewSchemaVersionHandler(app.SchemaVersions),
		SchemaFieldsHandler:  NewSchemaFieldsHandler(app.SchemaFields),
		ScopeHandler:         NewScopeHandler(app.Scope),
	}
}
