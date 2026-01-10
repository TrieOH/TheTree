package testing

import (
	"GoAuth/internal/adapters/http/router"
	"GoAuth/internal/database"
	"GoAuth/internal/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/gavv/httpexpect/v2"
	"github.com/spf13/viper"
)

// ============================================================================
// TEST FRAMEWORK - Core Infrastructure
// ============================================================================

// TestSuite manages the entire test environment
type TestSuite struct {
	Server *httptest.Server
	DB     *sql.DB
	t      *testing.T
}

func NewTestSuite(t *testing.T) *TestSuite {
	suite := &TestSuite{t: t}
	suite.setup()

	t.Cleanup(func() {
		suite.teardown()
	})

	return suite
}

func (s *TestSuite) setup() {
	viper.AutomaticEnv()

	if err := utils.LoadEd25519Keys(
		viper.GetString("JWT_PRIVATE_KEY"),
		viper.GetString("JWT_PUBLIC_KEY"),
	); err != nil {
		log.Fatal(err)
	}

	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024, // 10MB
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        "go-auth-test",
	})

	var err error
	s.DB, err = setupDatabase()
	if err != nil {
		s.t.Fatalf("DB setup failed: %v", err)
	}

	r := createTestRouter(s.DB)
	s.Server = httptest.NewServer(r)
}

func (s *TestSuite) teardown() {
	if s.Server != nil {
		s.Server.Close()
	}
	if s.DB != nil {
		_ = s.DB.Close()
	}
}

// NewClient creates a new API client for testing
func (s *TestSuite) NewClient(t *testing.T) *Client {
	return &Client{
		expect: httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  s.Server.URL,
			Reporter: httpexpect.NewAssertReporter(t),
		}),
		t: t,
	}
}

// ============================================================================
// HELPER FUNCTIONS - Keep these minimal
// ============================================================================

func setupDatabase() (*sql.DB, error) {
	db, err := database.WaitForDB(30 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}
	if err = database.RunMigrations(db, "../internal/database/migrations"); err != nil {
		return nil, fmt.Errorf("failed migrations: %w", err)
	}
	if err = database.SetJWTMasterKey(db); err != nil {
		return nil, err
	}
	return db, nil
}

func createTestRouter(db *sql.DB) http.Handler {
	return router.CreateTestRouter(db)
}
