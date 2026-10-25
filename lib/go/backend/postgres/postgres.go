// Package postgres is the kit module for a Backend's own Postgres database:
// the pool, the tx runner, and the constraint messages repos surface on
// violation. It never migrates: `just <svc>-goose` owns the schema.
package postgres

import (
	"context"

	"lib/backend"
	libconfig "lib/config"
	"lib/database"
)

// Module connects to the database described by cfg and fills Deps.DB and
// Deps.Tx. constraints maps the schema's constraint names to what a
// violation surfaces as (messages, or named errors per ADR-0016); boot fails
// when the schema has a constraint with no entry.
func Module(cfg libconfig.Postgres, constraints database.Constraints) backend.Module {
	return backend.Module{
		Name: "postgres",
		Open: func(_ context.Context, d *backend.Deps) (backend.Close, error) {
			pool, err := database.SetupDBWithoutMigrations(libconfig.DBConnectionConfig(cfg), constraints)
			if err != nil {
				return nil, err
			}
			d.DB = pool
			d.Tx = database.NewPGXTxRunner(pool)
			return func(context.Context) error {
				database.CloseDB(pool)
				return nil
			}, nil
		},
	}
}
