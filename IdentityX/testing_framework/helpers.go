package testing

import (
	"GoAuth/internal/adapters/http/router"
	"GoAuth/internal/database"
	"database/sql"
	"log"
	"net/http"
	"time"
)

// ============================================================================
// HELPER FUNCTIONS - Keep these minimal
// ============================================================================

func setupDatabase() (*sql.DB, error) {
	db, err := database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	if err := database.RunMigrations(db, "../internal/database/migrations"); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}
	if err := database.SetJWTMasterKey(db); err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func createTestRouter(db *sql.DB) http.Handler {
	return router.CreateTestRouter(db)
}
