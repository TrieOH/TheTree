package testing

import (
	"GoAuth/initialization"
	"GoAuth/internal/adapters/http/router"
	"GoAuth/internal/application"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// ============================================================================
// TEST FRAMEWORK - Core Infrastructure
// ============================================================================

// TestSuite manages the entire test environment
type TestSuite struct {
	Server *httptest.Server
	App    *application.Application
	DB     *pgxpool.Pool
	Redis  *redis.Client
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
	var goAuth initialization.GoauthApp

	initialization.LoadEnv(&goAuth)
	initialization.SetupFail()
	initialization.SetupFUN()
	initialization.SetupDB(&goAuth, "../internal/database/migrations")
	goAuth.Redis = initialization.SetupRedis(15 * time.Second)

	s.DB = goAuth.DB
	s.Redis = goAuth.Redis

	r, app := createTestRouter(s.DB, s.Redis)
	s.App = app
	s.Server = httptest.NewServer(r)
}

func (s *TestSuite) teardown() {
	if s.Server != nil {
		s.Server.Close()
	}
	if s.DB != nil {
		s.DB.Close()
	}
	if s.Redis != nil {
		s.Redis.Close()
	}
}

// NewClient creates a new API client for testing
func (s *TestSuite) NewClient(t *testing.T) *Client {
	return &Client{
		expect: httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  s.Server.URL,
			Reporter: httpexpect.NewAssertReporter(t),
		}),
		t:       t,
		baseURL: s.Server.URL,
	}
}

func createTestRouter(db *pgxpool.Pool, rdb *redis.Client) (http.Handler, *application.Application) {
	return router.CreateTestRouter(db, rdb)
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
