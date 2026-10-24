package main

import (
	"lib/backend"
	"lib/backend/identityx"
	"lib/backend/jobs"
	"lib/backend/postgres"
	"lib/errx"
	spec "payssage"
	"payssage/internal/app"
	"payssage/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	errx.Exit(backend.Run(backend.App{
		Name:   cfg.AppName,
		Spec:   spec.OpenAPISpec,
		Server: backend.ServerFrom(cfg.Server, cfg.CORS, cfg.ProfilePort),
		Modules: []backend.Module{
			postgres.Module(cfg.Postgres, app.ConstraintMessages()),
			identityx.Module(cfg.IdentityX, identityx.Debug(true), identityx.AcceptAPIKeys()),
			jobs.Module(cfg.JobsURL),
		},
		Build: app.Build(cfg),
	}), "payssage exited with error")
}
