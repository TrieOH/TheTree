package main

import (
	spec "Informd"
	"Informd/internal/app"
	"Informd/internal/config"
	"Informd/internal/problems"
	"lib/backend"
	"lib/backend/identityx"
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
			identityx.Module(cfg.IdentityX),
		},
		Build:    app.Build,
		Errors:   problems.Mapper(nil),
		Problems: problems.Docs,
	}), "informd exited with error")
}
