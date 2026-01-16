package application

import (
	"GoAuth/internal/adapters/email"
	"GoAuth/internal/adapters/persistence"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/application/authenticator"
	"GoAuth/internal/application/project"
	"GoAuth/internal/application/schema"
	"GoAuth/internal/application/schema_fields"
	"GoAuth/internal/application/schema_version"
	"GoAuth/internal/application/session"
	"GoAuth/internal/application/tokens"
	"GoAuth/internal/infrastructure"
	"GoAuth/internal/ports/inbounds"
)

type Application struct {
	Auth           inbounds.AuthService
	Project        inbounds.ProjectService
	Schema         inbounds.SchemaService
	SchemaVersions inbounds.SchemaVersionService
	SchemaFields   inbounds.SchemaFieldsService
	Session        inbounds.SessionService
	Authenticator  inbounds.RequestAuthenticator
}

func NewApplication(infra infrastructure.Infra) *Application {
	repos := persistence.NewRepositories(infra)

	mailBundle := email.NewBundle(infra)
	tokensBundle := tokens.NewBundle(repos)

	return &Application{
		Auth: auth.New(auth.Deps{
			Users:        repos.Users,
			Sessions:     repos.Sessions,
			Schemas:      repos.Schemas,
			Versions:     repos.SchemaVersions,
			Fields:       repos.SchemaFields,
			Projects:     repos.Projects,
			ProjectUsers: repos.ProjectUsers,
		}, infra, tokensBundle, mailBundle),
		Project: project.New(repos.Projects, infra.Tx),
		Schema: schema.New(schema.Deps{
			Schemas:  repos.Schemas,
			Versions: repos.SchemaVersions,
			Fields:   repos.SchemaFields,
			Projects: repos.Projects,
		}, infra.Tx),
		SchemaVersions: schema_version.New(schema_version.Deps{
			Schemas:  repos.Schemas,
			Versions: repos.SchemaVersions,
			Fields:   repos.SchemaFields,
			Projects: repos.Projects,
		}, infra.Tx),
		SchemaFields: schema_fields.New(schema_fields.Deps{
			Schemas:  repos.Schemas,
			Versions: repos.SchemaVersions,
			Fields:   repos.SchemaFields,
			Projects: repos.Projects,
		}, infra.Tx),
		Session: session.New(repos.Sessions, infra.Tx),
		Authenticator: authenticator.New(authenticator.Deps{
			Session:       repos.Sessions,
			TokenVerifier: tokensBundle.Verifier,
		}, infra.Tracer),
	}
}
