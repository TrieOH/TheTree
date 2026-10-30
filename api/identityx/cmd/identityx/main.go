package main

import (
	spec "IdentityX"
	"IdentityX/internal/app"
	"IdentityX/internal/config"
	"IdentityX/internal/problems"
	"lib/backend"
	"lib/backend/email"
	"lib/backend/jobs"
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
			email.Module(cfg.ToEmailConfig()),
			jobs.Module(cfg.JobsURL),
		},
		Build:    app.Build(cfg),
		Errors:   problems.Mapper(app.Problem),
		Problems: problems.Docs,
	}), "identityx exited with error")
}
