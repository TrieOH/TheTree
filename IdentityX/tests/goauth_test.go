package integration_test

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	database "GoAuth/internal/db"
	"GoAuth/internal/router"

	resp "github.com/MintzyG/GoResponse/response"
	"github.com/gavv/httpexpect/v2"
	"github.com/spf13/viper"
)

var Port string
var Db *sql.DB

func init() {
	viper.AutomaticEnv()

	Port = viper.GetString("PORT")
	if Port == "" {
		Port = "8080"
	}

	resp.SetConfig(resp.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024, // 10MB
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        "greet-test",
	})

	var err error
	Db, err = database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	if err := database.RunMigrations(Db, "../migrations"); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}
}

var serverUrl string

func runServer() {
	r := router.CreateTestRouter(Db)
	server := httptest.NewServer(r)
	serverUrl = server.URL
}

func createExpect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  serverUrl,
		Reporter: httpexpect.NewAssertReporter(t),
	})
}

func TestGreetService(t *testing.T) {
	runServer()
	defer Db.Close()

	t.Run("CreateUsers", func(t *testing.T) {
		runCreateUsers(t)
	})
}

type createUserContext struct {
	id string
}

func runCreateUsers(t *testing.T) {
	var ctx createUserContext

	completeUser := map[string]interface{}{
		"first_name": "Complete",
		"last_name":  "User",
	}

	t.Run("CreateCompleteUser", createUserSuccess(&ctx, completeUser))
	t.Run("GetCompleteUser", getUserAndVerify(&ctx, nil))
	t.Run("GetAndVerifyCompleteUser", getUserAndVerify(&ctx, &completeUser))

	incompleteUser := map[string]interface{}{
		"first_name": "IncompleteUser",
	}

	t.Run("CreateIncompleteUser", createUserSuccess(&ctx, incompleteUser))
	t.Run("GetIncompleteUser", getUserAndVerify(&ctx, nil))
	t.Run("GetAndVerifyIncompleteUser", getUserAndVerify(&ctx, &incompleteUser))
}

func getUserAndVerify(ctx *createUserContext, expected *map[string]interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		userObj := e.GET("/users/{id}", ctx.id).
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("data").Object()

		userObj.Value("id").NotNull()
		userObj.Value("first_name").NotNull()

		if expected != nil {
			for key, val := range *expected {
				userObj.Value(key).Equal(val)
			}
		}
	}
}

func createUserSuccess(ctx *createUserContext, user map[string]interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		r := e.POST("/users").
			WithHeader("Content-Type", "application/json").
			WithJSON(user).
			Expect().
			Status(http.StatusCreated).
			JSON().Object().Value("data").Object()

		for key, _ := range user {
			r.Value(key).NotNull()
		}

		ctx.id = r.Value("id").Raw().(string)
	}
}
