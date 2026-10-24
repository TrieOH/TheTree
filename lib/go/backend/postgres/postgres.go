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
// Deps.Tx. messages maps the schema's constraint names to client-facing
// messages; boot fails when the schema has a constraint with no message.
func Module(cfg libconfig.Postgres, messages database.ConstraintRegistry) backend.Module {
	return backend.Module{
		Name: "postgres",
		Open: func(_ context.Context, d *backend.Deps) (backend.Close, error) {
			pool, err := database.SetupDBWithoutMigrations(libconfig.DBConnectionConfig(cfg), messages)
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
