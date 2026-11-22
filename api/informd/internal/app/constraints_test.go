package app_test

import (
	"context"
	"testing"

	"Informd/internal/app"
	"lib/database"
	"lib/testdb"
)

// TestConstraintsMatchSchema — the registry and the migrated schema name
// the same constraints: a constraint missing here would answer a generic
// 500, and an entry naming nothing is one renamed or dropped.
func TestConstraintsMatchSchema(t *testing.T) {
	db := testdb.Postgres(t, "../../db/migrations")
	missing, unknown, err := database.RegistryDrift(context.Background(), db, app.Constraints())
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) > 0 {
		t.Errorf("constraints with no named error: %v", missing)
	}
	if len(unknown) > 0 {
		t.Errorf("entries naming no constraint in the schema: %v", unknown)
	}
}
