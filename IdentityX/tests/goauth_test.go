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
		DefaultModule:        "go-auth-test",
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

func TestGoAuthService(t *testing.T) {
	runServer()
	defer Db.Close()

	t.Run("RegisterTests", func(t *testing.T) {
		runRegisterTests(t)
	})
}

type createUserContext struct {
	id string
}

func runRegisterTests(t *testing.T) {
	// var ctx createUserContext
	success := map[string]interface{}{
		"email": "success@mail.com",
		"password":  "Str0ngP4ass!",
	}

  noEmail := map[string]interface{}{
		"email": "",
		"password":  "N0Email#S4d",
	}

  noPassword := map[string]interface{}{
		"email": "nopass@mail.com",
		"password":  "",
	}

	t.Run("RegisterSuccess", registerSuccess(&success))
	t.Run("RegisterNoEmail", registerNoEmail(&noEmail))
	t.Run("RegisterNoPassword", registerNoPassword(&noPassword))
}

// func getUserAndVerify(ctx *createUserContext, expected *map[string]interface{}) func(t *testing.T) {
// 	return func(t *testing.T) {
// 		e := createExpect(t)
//
// 		userObj := e.GET("/users/{id}", ctx.id).
// 			WithHeader("Content-Type", "application/json").
// 			Expect().
// 			Status(http.StatusOK).
// 			JSON().Object().Value("data").Object()
//
// 		userObj.Value("id").NotNull()
// 		userObj.Value("first_name").NotNull()
//
// 		if expected != nil {
// 			for key, val := range *expected {
// 				userObj.Value(key).Equal(val)
// 			}
// 		}
// 	}
// }

func registerSuccess(user *map[string]interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		r := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(user).
			Expect().
			Status(http.StatusCreated).
			JSON().Object().Value("message")

		r.String().Equal("Registered user")
	}
}

func registerNoEmail(user *map[string]interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(user).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("(email) is required")

		obj.Value("code").Number().Equal(400)
	}
}

func registerNoPassword(user *map[string]interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(user).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("(password) is required")

		obj.Value("code").Number().Equal(400)
	}
}
