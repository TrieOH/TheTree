package main

import (
	spec "Informd"
	"Informd/internal/app"
	"Informd/internal/config"
	"Informd/internal/problems"
	"lib/backend"
	"lib/backend/identityx"
	"lib/backend/postgres"
	"lib/errx"
)

func main() {
	cfg := config.LoadConfig()
	errx.Exit(backend.Run(backend.App{
		Name:   cfg.AppName,
		Spec:   spec.OpenAPISpec,
		Server: backend.ServerFrom(cfg.Server, cfg.CORS, cfg.ProfilePort),
		Modules: []backend.Module{
			postgres.Module(cfg.Postgres, app.Constraints()),
			identityx.Module(cfg.IdentityX, identityx.Debug(true)),
		},
		Build:    app.Build,
		Errors:   problems.Mapper(app.Problem),
		Problems: problems.Docs,
	}), "informd exited with error")
}
