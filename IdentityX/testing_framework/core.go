package testing

import (
	"GoAuth/internal/adapters/http/router"
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/crypto"
	"GoAuth/internal/database"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/gavv/httpexpect/v2"
	"github.com/oklog/ulid/v2"
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

	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024, // 10MB
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        "go-auth-test",
		ErrorHandler:         apierr.ErrToResp,
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

	queries := sqlc.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = queries.GetActiveSigningKeyForGoAuth(ctx)
	if err != nil {
		if apierr.IsNotFound(apierr.FromSQLC(err)) {
			// create new signing key
			pub, priv, err := crypto.GenerateEd25519()
			if err != nil {
				log.Fatalf("failed to generate GoAuth key: %v", err)
			}
			defer zero(priv)

			kid := "goauth:" + ulid.Make().String()
			expiresAt := time.Now().Add(90 * 24 * time.Hour)

			_, err = queries.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
				Kid:        kid,
				ProjectID:  nil,
				KeyType:    "goauth",
				Algorithm:  "EdDSA",
				PublicKey:  pub,
				PrivateKey: priv,
				Usage:      "sign",
				Status:     "active",
				ExpiresAt:  expiresAt,
			})

			if err != nil {
				// rely on DB uniqueness as safety net
				if apierr.IsUniqueViolation(err) {
					log.Println("GoAuth signing key already created by another instance")
				} else {
					log.Fatalf("failed to create GoAuth signing key: %v", err)
				}
			} else {
				log.Println("Created GoAuth signing key")
			}
		} else {
			log.Fatalf("failed checking GoAuth signing key: %v", err)
		}
	}
	return db, nil
}

func createTestRouter(db *sql.DB) http.Handler {
	return router.CreateTestRouter(db)
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
