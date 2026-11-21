package main

import (
	spec "IdentityX"
	"IdentityX/internal/app"
	"IdentityX/internal/app/errmap"
	"IdentityX/internal/config"
	"IdentityX/internal/problems"
	"lib/backend"
	"lib/backend/email"
	"lib/backend/jobs"
	"lib/backend/postgres"
	libconfig "lib/config"
	"lib/errx"
)

func main() {
	cfg := libconfig.Load[config.Config]()
	errx.Exit(backend.Run(backend.App{
		Name:   cfg.AppName,
		Spec:   spec.OpenAPISpec,
		Server: cfg.Server,
		CORS:   cfg.CORS,
		Modules: []backend.Module{
			postgres.Module(cfg.Postgres, app.Constraints()),
			email.Module(cfg.SMTP),
			jobs.Module(cfg.Jobs),
		},
		Build:    app.Build(cfg),
		Errors:   problems.Mapper(errmap.Problem),
		Problems: problems.Docs,
	}), "identityx exited with error")
}
