package main

import (
	"lib/backend"
	"lib/backend/email"
	"lib/backend/identityx"
	"lib/backend/jobs"
	"lib/backend/postgres"
	libconfig "lib/config"
	"lib/errx"
	spec "univents"
	"univents/internal/app"
	"univents/internal/app/errmap"
	"univents/internal/config"
	"univents/internal/problems"
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
			identityx.Module(cfg.IdentityX, identityx.RequireVerifiedEmail(cfg.RequireVerifiedEmail)),
			email.Module(cfg.SMTP),
			jobs.Module(cfg.Jobs),
		},
		Build:    app.Build(cfg),
		Errors:   problems.Mapper(errmap.Problem),
		Problems: problems.Docs,
	}), "univents exited with error")
}
