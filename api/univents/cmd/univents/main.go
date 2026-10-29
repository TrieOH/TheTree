package main

import (
	"lib/backend"
	"lib/backend/email"
	"lib/backend/identityx"
	"lib/backend/jobs"
	"lib/backend/postgres"
	"lib/errx"
	spec "univents"
	"univents/internal/app"
	"univents/internal/app/errmap"
	"univents/internal/config"
)

func main() {
	cfg := config.Load()
	errx.Exit(backend.Run(backend.App{
		Name:   cfg.AppName,
		Spec:   spec.OpenAPISpec,
		Server: backend.ServerFrom(cfg.Server, cfg.CORS, cfg.ProfilePort),
		Modules: []backend.Module{
			postgres.Module(cfg.Postgres, app.Constraints()),
			identityx.Module(cfg.IdentityX,
				identityx.Debug(cfg.DebugMode),
				identityx.RequireVerifiedEmail(cfg.RequireVerifiedEmail)),
			email.Module(cfg.ToEmailConfig()),
			jobs.Module(cfg.JobsURL),
		},
		Build:  app.Build(cfg),
		Errors: errmap.Errors(),
	}), "univents exited with error")
}
