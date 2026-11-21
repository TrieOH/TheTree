package main

import (
	"lib/backend"
	"lib/backend/identityx"
	"lib/backend/jobs"
	"lib/backend/postgres"
	libconfig "lib/config"
	"lib/errx"
	spec "payssage"
	"payssage/internal/app"
	"payssage/internal/app/errmap"
	"payssage/internal/config"
	"payssage/internal/problems"
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
			identityx.Module(cfg.IdentityX, identityx.AcceptAPIKeys()),
			jobs.Module(cfg.Jobs),
		},
		Build:    app.Build(cfg),
		Errors:   problems.Mapper(errmap.Problem),
		Problems: problems.Docs,
	}), "payssage exited with error")
}
