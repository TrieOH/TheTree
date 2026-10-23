package app

import (
	"context"

	"Informd/internal/authz"
	"Informd/internal/handlers"
	"Informd/internal/repos"
	"Informd/internal/services"
	"Informd/internal/sqlc"
	"lib/backend"

	"github.com/go-chi/chi/v5"
)

// Build wires Informd on the kit (ADR-0014): repos and operations over the
// postgres module, auth from the identityx module, and the strict server.
func Build(_ context.Context, d *backend.Deps, r chi.Router) error {
	repos := repos.New(sqlc.New(d.DB))
	ops := services.NewOperations(repos, authz.New(repos.Forms, repos.Namespaces), d.Tx)

	chains, err := d.Chains(d.IDX.Auth)
	if err != nil {
		return err
	}
	mountStrict(r, handlers.NewServer(ops), chains)
	return nil
}
